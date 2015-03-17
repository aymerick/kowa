package server

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"sync"
	"time"

	"github.com/aymerick/kowa/builder"
	"github.com/aymerick/kowa/models"
	"github.com/spf13/viper"
)

const (
	WORKERS_NB        = 10
	JOBS_QUEUE_LEN    = 100
	WORKERS_QUEUE_LEN = 100
)

type BuildMaster struct {
	workers   []*BuildWorker
	workersWG *sync.WaitGroup

	workersChan chan *BuildJob
	jobsChan    chan *BuildJob
	resultsChan chan *BuildJob

	currentJobs   map[string]*BuildJob
	throttledJobs map[string]*BuildJob

	stopChan chan bool
}

type BuildJob struct {
	siteId string
	failed bool
}

type BuildWorker struct {
	id int

	inputChan  chan *BuildJob
	outputChan chan *BuildJob
	stopChan   chan bool
}

//
// BuildMaster
//

func NewBuildMaster() *BuildMaster {
	result := &BuildMaster{
		workers:   make([]*BuildWorker, WORKERS_NB),
		workersWG: &sync.WaitGroup{},

		workersChan: make(chan *BuildJob, WORKERS_QUEUE_LEN),
		jobsChan:    make(chan *BuildJob, JOBS_QUEUE_LEN),
		resultsChan: make(chan *BuildJob, WORKERS_QUEUE_LEN),

		currentJobs:   make(map[string]*BuildJob),
		throttledJobs: make(map[string]*BuildJob),
	}

	return result
}

func (master *BuildMaster) NewBuildJob(siteId string) *BuildJob {
	return &BuildJob{
		siteId: siteId,
	}
}

func (master *BuildMaster) NewBuildWorker(workerId int) *BuildWorker {
	return &BuildWorker{
		id: workerId,

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
					log.Printf("[build] Job %s throttled", jobKey)

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

				log.Printf("[build] Job %s done", jobKey)

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
	dir := path.Join(viper.GetString("working_dir"), viper.GetString("output_dir"))
	port := viper.GetInt("serve_output_port")

	log.Println("[build] Serving built sites on port:", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), http.FileServer(http.Dir(dir))))
}

// Initialize and start all workers
func (master *BuildMaster) startWorkers() {
	for i := 0; i < WORKERS_NB; i++ {
		master.workers[i] = master.NewBuildWorker(i)
		master.workers[i].run(master.workersWG)
	}

	log.Printf("[build] Started %d workers", WORKERS_NB)
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

func (master *BuildMaster) enqueueJob(job *BuildJob) {
	master.jobsChan <- job

	log.Printf("[build] Job %s enqueued", job.key())
}

func (master *BuildMaster) launchSiteBuild(site *models.Site) {
	master.enqueueJob(master.NewBuildJob(site.Id))
}

//
// BuildJob
//

// Computes job uniq key
func (job *BuildJob) key() string {
	return job.siteId
}

//
// BuildWorker
//

// Starts worker
func (worker *BuildWorker) run(wg *sync.WaitGroup) {
	if worker.stopChan != nil {
		// already running
		return
	}

	// run
	go func() {
		// add oursel to workgroup and release on exit
		wg.Add(1)
		defer wg.Done()

		// setup stop channel
		worker.stopChan = make(chan bool)
		defer close(worker.stopChan)

		ended := false

		for !ended {
			select {
			case job := <-worker.inputChan:
				log.Printf("[build] Job %s taken by worker %d", job.key(), worker.id)

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
func (worker *BuildWorker) executeJob(job *BuildJob) {
	// get site
	site := models.NewDBSession().FindSite(job.siteId)
	if site == nil {
		log.Printf("[build] Job %s failed with worker %d: site not found", job.key(), worker.id)

		job.failed = true
		return
	}

	// builder config
	config := &builder.SiteBuilderConfig{
		WorkingDir: viper.GetString("working_dir"),
		OutputDir:  path.Join(viper.GetString("output_dir"), site.Id),
		BasePath:   path.Join("/", site.Id),
	}

	builder := builder.NewSiteBuilder(site, config)

	if builder.Build(); builder.HaveError() {
		// job failed
		job.failed = true
	} else {
		// update BuiltAt anchor
		site.SetBuiltAt(time.Now())
	}
}

// Stop worker
func (worker *BuildWorker) stop() {
	// wait for worker to stop
	worker.stopChan <- true
	<-worker.stopChan
}
