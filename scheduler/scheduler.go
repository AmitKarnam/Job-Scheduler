package scheduler

import (
	"sync"
	"time"

	"github.com/AmitKarnam/Job-Scheduler/heap"
)

const (
	// TimeWindow describes the time window for which the jobs will be loaded from db into memory
	TimeWindow = 5 * time.Minute
	// LookAheadTimeWindow describes the time interval in which we load and populate the heap. This should be done ahead of TimeWindow to ensure that jobs to be executed are not exhausted.
	LookAheadTimeWindow = 3 * time.Minute
)

type Scheduler interface {
	buildHeap(readLevel, ackLevel string) heap.Heap
	monitorAckLevelandReadLevel()
	jobExecutor()
	Start()
	Stop()
}

type scheduler struct {
	ReadLevel  time.Time
	AckLevel   time.Time
	jobMinHeap heap.Heap
	wg         *sync.WaitGroup
	stopChan   chan bool
}

func InitialiseScheduler(timeWindow int) Scheduler {
	// All the following logic is one time initialisation, once the scheduler is intialised, we use it's attributes to make the updates
	// read Acklevel from database
	// read Acklevel from database
	ackLevel, readLevel := loadAckLevelAndReadLevel()
	// return the scheduler with empty jobMinheap, then when the scheduler is started populate it.
	return &scheduler{
		ReadLevel:  readLevel,
		AckLevel:   ackLevel,
		jobMinHeap: heap.NewHeap(),
	}
}

// Currently a dummy method that provides the AckLevel and ReadLevel.
// function used to fetch the current AckLevel and ReadLevel stored in the database.
func loadAckLevelAndReadLevel() (time.Time, time.Time) {
	return time.Now().UTC(), time.Now().UTC()
}

// buildHeap is a initialiser method that help to build the initial in memory min-heap based on teh execution window ( 'N' ) time window
func (s *scheduler) buildHeap(readLevel, ackLevel string) heap.Heap {
	return heap.Heap
}

func (s *scheduler) monitorAckLevelandReadLevel() {
	// Get current time and ReadLevel
	// If current time is 2-minutes less than ReadLevel; Populate the heap with data from ReadLevel to next 'N' minutes
	// What's the trigger?? => How will know when to fire
	defer s.wg.Done()

	ticker := time.NewTicker(LookAheadTimeWindow)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopChan:
			return
		case <-ticker.C:
			// Load jobs from DB between ReadLevel and ReadLevel + window
			// Insert into heap
			// Update ReadLevel in DB
		}
	}

}

func (s *scheduler) jobExecutor() {
	defer s.wg.Done()

	for {
		jobToExecute := s.jobMinHeap.Peak()

		select {
		//TODO: Case if a new job is added to the heap with an earlier execution time than the current job at the top of the heap, how to handle that? - We can add a new channel to the scheduler struct called 'heapUpdateChan' and whenever a new job is added to the heap, we can send a signal to that channel. In the jobExecutor, we can listen to that channel and if we receive a signal, we can re-evaluate the top of the heap and adjust the sleep time accordingly.
		case <-s.stopChan:
			return
		case <-time.After(time.Until(jobToExecute.NextExecutionTime)):
			// Execute the job
			err := jobToExecute.Execute()
			if err != nil {
				// Handle error (e.g., log it, retry logic, etc.)
			}
			s.jobMinHeap.DeleteMin() // Remove the executed job from the heap
			//TODO: If the job is recurring, compute its next execution time and update the DB entry for that job; It will be automatically be picked up by the monitorAckLevelandReadLevel and added to the heap when it's next_execution_time is within the ReadLevel and ReadLevel + window
			return
		}
	}

}

func (s *scheduler) Start() {
	// Use ack level and read level to create a min-heap
	// Case 1: The ReadLevel and AckLevel are empty ( new instance of job scheduler ): Set both of them to current timestamp; load the min-heap with the jobs that execute within next 'N' time units
	if s.AckLevel.IsZero() && s.ReadLevel.IsZero() {

	}
	// Case 2: ReadLevel and AckLevel are far in the past from current time ( Job scheduler crash or stopped ): Load all the jobs from the AckLevel to the current timestamp ( should thier execution be taken care by a seperate worker? ), Load all the jobs from current timestamp + 'N' time units
	if s.AckLevel < time.Now().Format() {
		// Load all the jobs in the from AckLevel to current time.
		// Async: Start the backlog job worker to execute the jobs in backlog
		// Set current time as AckLevel and ReadLevel; Start loadin jobs from current time to next 'N' minutes
	}
	// Case 3: ReadLevel or AckLevel are in the future; Alert for a drifted system clock; Ask users to sync clock; Provide instructions; Exit
	// Move the AckLevel and ReadLevel after the above step to their correct timestamp => Flush to DB
	// Start jobExecutor as a go-routine
	// Start monitorAckLevelandReadLevel as a go routine
	s.buildHeap(s.ReadLevel.Format("2006-01-02T15:04:05Z"), s.AckLevel.Format("2006-01-02T15:04:05Z"))
	s.wg.Add(2)
	go s.jobExecutor()
	go s.monitorAckLevelandReadLevel()
	s.wg.Wait()

	return

}

func (s *scheduler) Stop() {
	// Need to understand more on the behaviour and handle accordingly
	close(s.stopChan)
}
