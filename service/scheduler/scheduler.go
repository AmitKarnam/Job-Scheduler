package scheduler

import (
	"fmt"
	"sync"
	"time"

	"github.com/AmitKarnam/Job-Scheduler/heap"
	"github.com/AmitKarnam/Job-Scheduler/models"
	"github.com/AmitKarnam/Job-Scheduler/repository"
	"github.com/AmitKarnam/Job-Scheduler/worker/backlog"
)

const (
	// TimeWindow describes the time window for which the jobs will be loaded from db into memory
	TimeWindow = 5 * time.Minute
	// LookAheadTimeWindow describes the time interval in which we load and populate the heap. This should be done ahead of TimeWindow to ensure that jobs to be executed are not exhausted.
	LookAheadTimeWindow = 3 * time.Minute
)

type SchedulerService interface {
	Start()
	Stop()
}

type scheduler struct {
	readLevel  time.Time
	ackLevel   time.Time
	jobMinHeap heap.Heap
	jobRepo    repository.JobRepo
	wg         *sync.WaitGroup
	stopChan   chan bool
}

func InitialiseScheduler(jobRepo *repository.JobRepo) SchedulerService {
	// All the following logic is one time initialisation, once the scheduler is intialised, we use it's attributes to make the updates
	// read Acklevel from database
	// read Acklevel from database
	ackLevel, readLevel := loadAckLevelAndReadLevel()
	// return the scheduler with empty jobMinheap, then when the scheduler is started populate it.
	return &scheduler{
		readLevel:  readLevel,
		ackLevel:   ackLevel,
		jobMinHeap: heap.NewHeap(),
		jobRepo:    *jobRepo,
	}
}

// Currently a dummy method that provides the AckLevel and ReadLevel.
// function used to fetch the current AckLevel and ReadLevel stored in the database.
func loadAckLevelAndReadLevel() (time.Time, time.Time) {
	return time.Now().UTC(), time.Now().UTC()
}

// buildHeap is a initialiser method that help to build the initial in memory min-heap based on teh execution window ( 'N' ) time window
func (s *scheduler) buildHeap(readLevel, ackLevel time.Time) heap.Heap {
	// Build heap based on the readLevel and ackLevel; Fetch the jobs from database that are scheduled to be executed between readLevel and readLevel + TimeWindow; Insert those jobs into the heap; Return the heap
	return heap.NewHeap()
}

func (s *scheduler) monitorAckLevelandReadLevel() {
	defer s.wg.Done()

	ticker := time.NewTicker(LookAheadTimeWindow)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopChan:
			return
		case <-ticker.C:
			// Load jobs from DB between ReadLevel and ReadLevel + window
			jobsToBeAdded := fetchJobsFromDB(s.readLevel, s.readLevel.Add(TimeWindow))
			// Insert into heap
			for _, job := range jobsToBeAdded {
				s.jobMinHeap.Insert(job)
			}
			// Update ReadLevel in DB
			s.readLevel = s.readLevel.Add(TimeWindow)
		}
	}

}

func (s *scheduler) jobExecutor() {
	defer s.wg.Done()

	for {
		jobToExecute := s.jobMinHeap.Peak()

		// handling empty job
		if jobToExecute.NextExecutionTime.IsZero() {
			select {
			case <-s.stopChan:
				return
			case <-time.After(1 * time.Second):
				// periodic recheck if any job was added
				continue
			}
		}

		select {
		//TODO: Case if a new job is added to the heap with an earlier execution time than the current job at the top of the heap, how to handle that? - We can add a new channel to the scheduler struct called 'heapUpdateChan' and whenever a new job is added to the heap, we can send a signal to that channel. In the jobExecutor, we can listen to that channel and if we receive a signal, we can re-evaluate the top of the heap and adjust the sleep time accordingly.
		case <-s.stopChan:
			return
		// TODO: This flow of job execution should be async in nature
		case <-time.After(time.Until(jobToExecute.NextExecutionTime)):
			// Execute the job
			err := jobToExecute.Execute()
			if err != nil {
				// Handle error (e.g., log it, retry logic, etc.)
				fmt.Printf("Error executing job: %v\n", err)
			}
			s.jobMinHeap.DeleteMin() // Remove the executed job from the heap
			//TODO: If the job is recurring, compute its next execution time and update the DB entry for that job; It will be automatically be picked up by the monitorAckLevelandReadLevel and added to the heap when it's next_execution_time is within the ReadLevel and ReadLevel + window
			return
		}
	}

}

func (s *scheduler) Start() {
	// Case 1: ReadLevel and AckLevel are far in the past from current time ( Job scheduler crash or stopped ): Load all the jobs from the AckLevel to the current timestamp ( should thier execution be taken care by a seperate worker ), Load all the jobs from current timestamp + 'N' time units
	if s.ackLevel.Before(time.Now()) && s.readLevel.Before(time.Now()) {
		// Load all the jobs in the from AckLevel to current time.
		// Async: Start the backlog job worker to execute the jobs in backlog
		// Set current time as AckLevel and ReadLevel; Start loadin jobs from current time to next 'N' minutes

		// List of backlog jobs
		backlogJobList := []models.Job{}

		backlogJobWorker := backlog.NewJobBacklogWorker()

		backlogJobWorker.Process()

		for _, backlogJob := range backlogJobList {
			backlogJobWorker.Add(backlogJob)
		}

		backlogJobWorker.Stop()

	}
	// Case 2: ReadLevel or AckLevel are in the future; Alert for a drifted system clock; Ask users to sync clock; Provide instructions; Exit
	if s.ackLevel.After(time.Now()) || s.readLevel.After(time.Now()) {
		fmt.Println("System clock is drifted. Please sync your clock with an NTP server and restart the scheduler.")
		return
	}

	// Case 3: Fresh job scheduler instance with no jobs in the database: Both ReadLevel and AckLevel are at the current timestamp; Start loading jobs from current time to next 'N' minutes
	s.readLevel = time.Now().UTC()
	s.ackLevel = time.Now().UTC()
	// Move the AckLevel and ReadLevel after the above step to their correct timestamp => Flush to DB
	// Start jobExecutor as a go-routine
	// Start monitorAckLevelandReadLevel as a go routine
	s.buildHeap(s.readLevel, s.readLevel.Add(TimeWindow))
	s.wg.Add(2)
	go s.jobExecutor()
	go s.monitorAckLevelandReadLevel()
	s.wg.Wait()
}

func (s *scheduler) Stop() {
	// Need to understand more on the behaviour and handle accordingly
	close(s.stopChan)
}

func fetchJobsFromDB(readLevel time.Time, maxTime time.Time) []models.Job {
	// Ideally run a DB query with next_execution_time as index and fetch all the jobs within the window readLevela nd maxTime
	return []models.Job{}
}
