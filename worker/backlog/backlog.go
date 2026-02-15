package backlog

import (
	"log"
	"sync"

	"github.com/AmitKarnam/Job-Scheduler/models"
)

// Backlog is responsible for managing the queue of jobs that are pending execution. It provides an interface for adding new jobs, retrieving jobs for execution, and updating the status of jobs after execution. The backlog ensures that jobs are executed in a timely manner and handles any necessary retries or error handling.
var (
	BacklogProcessingWorkerCount = 5
	MaxbacklogJob                = 100
)

type JobsBacklog interface {
	Add(models.Job)
	Process()
	Stop()
}

type backlog struct {
	queue   chan models.Job
	workers int
	wg      sync.WaitGroup
}

func NewJobBacklogWorker() JobsBacklog {
	return &backlog{
		queue:   make(chan models.Job, 100),
		workers: BacklogProcessingWorkerCount,
	}
}

func (bl *backlog) Add(job models.Job) {
	bl.queue <- job
}

func (bl *backlog) Process() {
	for i := 0; i < bl.workers; i++ {
		bl.wg.Add(1)
		go bl.worker(i)
	}
}

func (bl *backlog) Stop() {
	// close the queue to signal workers to exit
	close(bl.queue)
	bl.wg.Wait()
}

func (bl *backlog) worker(id int) {
	defer bl.wg.Done()

	for job := range bl.queue {
		// Execute job
		// Set it's next execution time based on the condition of the job
		if err := job.Execute(); err != nil {
			log.Printf("backlog worker %d: job execute error: %v", id, err)
		}
	}
}
