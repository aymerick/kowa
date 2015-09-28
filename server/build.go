package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"sync"
	"time"

	"github.com/aymerick/kowa/builder"
	"github.com/aymerick/kowa/models"
	"github.com/spf13/viper"
)

const (
	workersNb       = 10
	jobsQueueLen    = 100
	workersQueueLen = 100

	// job kinds
	jobKindBuild  = "build"
	jobKindDelete = "delete"
)

// BuildMaster handles build workers
type BuildMaster struct {
	workers   []*buildWorker
	workersWG *sync.WaitGroup

	workersChan chan *buildJob
	jobsChan    chan *buildJob
	resultsChan chan *buildJob

	currentJobs   map[string]*buildJob
	throttledJobs map[string]*buildJob

	stopChan chan bool
}

type buildJob struct {
	kind     string
	siteID   string
	buildDir string
	failed   bool
}

type buildWorker struct {
	id int

	inputChan  chan *buildJob
	outputChan chan *buildJob
	stopChan   chan bool
}

//
// BuildMaster
//

// NewBuildMaster instanciates a new build master
func NewBuildMaster() *BuildMaster {
	result := &BuildMaster{
		workers:   make([]*buildWorker, workersNb),
		workersWG: &sync.WaitGroup{},

		workersChan: make(chan *buildJob, workersQueueLen),
		jobsChan:    make(chan *buildJob, jobsQueueLen),
		resultsChan: make(chan *buildJob, workersQueueLen),

		currentJobs:   make(map[string]*buildJob),
		throttledJobs: make(map[string]*buildJob),
	}

	return result
}

func (master *BuildMaster) newBuildJob(kind string, siteID string, buildDir string) *buildJob {
	return &buildJob{
		kind:     kind,
		siteID:   siteID,
		buildDir: buildDir,
	}
}

func (master *BuildMaster) newBuildWorker(workerID int) *buildWorker {
	return &buildWorker{
		id: workerID,

		inputChan:  master.workersChan,
		outputChan: master.resultsChan,
	}
}

// Starts build master
func (master *BuildMaster) run() {
	if master.stopChan != nil {
		// master already running
		return
	}

	// run master
	go func() {
		// setup stop channel
		master.stopChan = make(chan bool)
		defer close(master.stopChan)

		// start workers
		master.startWorkers()
		defer master.stopWorkers()

		ended := false

		for !ended {
			select {
			case job := <-master.jobsChan:
				// new build job received
				jobKey := job.key()

				if master.currentJobs[jobKey] != nil {
					log.Printf("[build] %s job %s throttled", job.kind, jobKey)

					// a worker is already processing that job
					master.throttledJobs[jobKey] = job
				} else {
					master.currentJobs[jobKey] = job

					// dispatch to workers
					master.workersChan <- job
				}

			case job := <-master.resultsChan:
				// build job ended
				jobKey := job.key()

				log.Printf("[build] %s job %s done", job.kind, jobKey)

				// remove from current jobs
				delete(master.currentJobs, jobKey)

				if newJob := master.throttledJobs[jobKey]; newJob != nil {
					delete(master.throttledJobs, jobKey)

					// enqueue throttled job
					master.enqueueJob(newJob)
				}

			case <-master.stopChan:
				ended = true
			}
		}

		log.Printf("[build] Master is shutdowning")
	}()

	log.Printf("[build] Master launched")
}

// Serve built sites
func (master *BuildMaster) serveSites() {
	dir := viper.GetString("output_dir")
	port := viper.GetInt("serve_output_port")

	log.Println("[build] Serving built sites on port:", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), http.FileServer(http.Dir(dir))))
}

// Initialize and start all workers
func (master *BuildMaster) startWorkers() {
	for i := 0; i < workersNb; i++ {
		master.workers[i] = master.newBuildWorker(i)
		master.workers[i].run(master.workersWG)
	}

	log.Printf("[build] Started %d workers", workersNb)
}

// Stop all workers
func (master *BuildMaster) stopWorkers() {
	for _, worker := range master.workers {
		go worker.stop()
	}

	log.Printf("[build] Waiting for all workers to stop")

	// wait for all workers to stop
	master.workersWG.Wait()

	log.Printf("[build] All workers stopped")

	master.workers = nil
}

// Stops build master
func (master *BuildMaster) stop() {
	if master.stopChan == nil {
		// master is not running
		return
	}

	// wait for master to stop
	master.stopChan <- true
	<-master.stopChan
	master.stopChan = nil

	close(master.jobsChan)
	close(master.workersChan)
}

func (master *BuildMaster) enqueueJob(job *buildJob) {
	master.jobsChan <- job

	log.Printf("[build] Job %s enqueued", job.key())
}

func (master *BuildMaster) launchSiteBuild(site *models.Site) {
	master.enqueueJob(master.newBuildJob(jobKindBuild, site.ID, ""))
}

func (master *BuildMaster) launchSiteDeletion(site *models.Site, buildDir string) {
	master.enqueueJob(master.newBuildJob(jobKindDelete, site.ID, buildDir))
}

//
// BuildJob
//

// Computes job uniq key
func (job *buildJob) key() string {
	return job.siteID
}

//
// BuildWorker
//

// Starts worker
func (worker *buildWorker) run(wg *sync.WaitGroup) {
	if worker.stopChan != nil {
		// already running
		return
	}

	// run
	go func() {
		// add ourself to workgroup and release on exit
		wg.Add(1)
		defer wg.Done()

		// setup stop channel
		worker.stopChan = make(chan bool)
		defer close(worker.stopChan)

		ended := false

		for !ended {
			select {
			case job := <-worker.inputChan:
				log.Printf("[build] %s job %s taken by worker %d", job.kind, job.key(), worker.id)

				// execute job
				worker.executeJob(job)

				// @todo Benchmark execution time
				// @todo Handle job failure
				// @todo Rescue job crash

				// send result
				worker.outputChan <- job

			case <-worker.stopChan:
				ended = true
			}
		}
	}()
}

// Execute Job
func (worker *buildWorker) executeJob(job *buildJob) {
	switch job.kind {
	case jobKindBuild:
		worker.buildSite(job)
	case jobKindDelete:
		worker.deleteSite(job)
	default:
		panic("wat")
	}
}

func (worker *buildWorker) buildSite(job *buildJob) {
	// get site
	site := models.NewDBSession().FindSite(job.siteID)
	if site == nil {
		log.Printf("[build] %s job %s failed with worker %d: site not found", job.kind, job.key(), worker.id)

		job.failed = true
		return
	}

	builder := builder.NewSiteBuilder(site)

	if builder.Build(); builder.HaveError() {
		// job failed
		job.failed = true

		builder.DumpErrors()
	} else {
		// update BuiltAt anchor
		site.SetBuiltAt(time.Now())
	}
}

func (worker *buildWorker) deleteSite(job *buildJob) {
	dirPath := path.Join(viper.GetString("output_dir"), job.buildDir)
	if _, err := os.Stat(dirPath); !os.IsNotExist(err) {
		if errRem := os.RemoveAll(dirPath); errRem != nil {
			job.failed = true
		}
	}
}

// Stop worker
func (worker *buildWorker) stop() {
	// wait for worker to stop
	worker.stopChan <- true
	<-worker.stopChan
}
