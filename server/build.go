package server

import (
	"log"
	"sync"
	"time"

	"github.com/aymerick/kowa/models"
)

const (
	WORKERS_NB        = 10
	JOBS_QUEUE_LEN    = 100
	WORKERS_QUEUE_LEN = 100
)

type BuildMaster struct {
	stopChan chan bool

	jobsChan    chan *BuildJob
	workersChan chan *BuildJob
	jobDoneChan chan *BuildJob

	currentJobs   map[string]*BuildJob
	throttledJobs map[string]*BuildJob

	workers []*BuildWorker
}

type BuildJob struct {
	siteId string
}

type BuildWorker struct {
	id           int
	incomingChan chan *BuildJob
	doneChan     chan *BuildJob
	stopChan     chan bool
	running      bool
}

//
// BuildMaster
//

func NewBuildMaster() *BuildMaster {
	result := &BuildMaster{
		jobsChan:    make(chan *BuildJob, JOBS_QUEUE_LEN),
		workersChan: make(chan *BuildJob, WORKERS_QUEUE_LEN),
		jobDoneChan: make(chan *BuildJob, WORKERS_QUEUE_LEN),

		currentJobs:   make(map[string]*BuildJob),
		throttledJobs: make(map[string]*BuildJob),

		workers: make([]*BuildWorker, WORKERS_NB),
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
		id:           workerId,
		incomingChan: master.workersChan,
		doneChan:     master.jobDoneChan,
		stopChan:     make(chan bool),
	}
}

// Starts build master
func (master *BuildMaster) run() {
	if master.stopChan != nil {
		// master already running
		return
	}

	master.stopChan = make(chan bool)

	// start workers
	master.startWorkers()

	// run master
	go func() {
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

			case job := <-master.jobDoneChan:
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
				master.stopWorkers()
				close(master.stopChan)
				ended = true
			}
		}

		log.Printf("[build] Master ended")
	}()

	log.Printf("[build] Master started")
}

// Initialize and start all workers
func (master *BuildMaster) startWorkers() {
	for i := 0; i < WORKERS_NB; i++ {
		master.workers[i] = master.NewBuildWorker(i)
		master.workers[i].run()
	}
}

// Stop all workers
func (master *BuildMaster) stopWorkers() {
	wg := &sync.WaitGroup{}
	wg.Add(WORKERS_NB)

	for _, worker := range master.workers {
		go worker.stop(wg)
	}

	// wait for all workers to stop
	wg.Wait()

	master.workers = nil
}

// Stops build master
// @todo Use that method !
func (master *BuildMaster) stop() {
	if master.stopChan == nil {
		// master is not running
		return
	}

	// wait for master to stop
	master.stopChan <- true
	<-master.stopChan

	close(master.stopChan)
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
func (worker *BuildWorker) run() {
	if worker.running {
		// already running
		return
	}

	// run
	go func() {
		ended := false

		for !ended {
			select {
			case job := <-worker.incomingChan:
				log.Printf("[build] Job %s taken by worker %d", job.key(), worker.id)

				// @todo Executes job
				time.Sleep(5 * time.Second)

				worker.doneChan <- job

			case <-worker.stopChan:
				close(worker.stopChan)
				ended = true
			}
		}

		log.Printf("[build] Worker %d ended", worker.id)
	}()

	log.Printf("[build] Worker %d started", worker.id)

	worker.running = true
}

// Stop worker
func (worker *BuildWorker) stop(wg *sync.WaitGroup) {
	defer wg.Done()

	// wait for worker to stop
	worker.stopChan <- true
	<-worker.stopChan
}
