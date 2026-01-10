package scheduler

import (
	"heap"
	"sync"
	"time"
)

const ()

type Scheduler interface {
	buildHeap() heap.heap
	monitorAckLevelandReadLevel()
	jobExecutor()
	Start() error
	Stop()
}

type scheduler struct {
	ReadLevel  string
	AckLevel   string
	jobMinHeap heap.heap
	wg         *sync.WaitGroup
	stopChan   chan bool
}

func InitialiseScheduler() Scheduler {
	// All the following logic is one time initialisation, once the scheduler is intialised, we use it's attributes to make the updates
	// read Acklevel from database
	// read Acklevel from database
	ackLevel, readLevel := loadAckLevelAndReadLevel()
	// return the scheduler with empty jobMinheap, then when the scheduler is started populate it.
	return &scheduler{
		ReadLevel:  readLevel,
		AckLevel:   ackLevel,
		jobMinHeap: heap.Heap{},
	}
}

func loadAckLevelAndReadLevel() (string, string) {
	return time.Now().Format(time.RFC3339), time.Now().Format(time.RFC3339)
}

func (s *scheduler) buildHeap() heap.Heap {
	return heap.Heap
}

func (s *scheduler) monitorAckLevelandReadLevel() {}

func (s *scheduler) jobExecutor() {}

func (s *scheduler) Start() error {
	// Use ack level and read level to create a min-heap
	// Case 1: The ReadLevel and AckLevel are empty ( new instance of job scheduler ): Set both of them to current timestamp; load the min-heap with the jobs that execute within next 'N' time units
	// Case 2: ReadLevel and AckLevel are far in the past from current time ( Job scheduler crash or stopped ): Load all the jobs from the AckLevel to the current timestamp ( should thier execution be taken care by a seperate worker? ), Load all the jobs from current timestamp + 'N' time units
	// Case 3: ReadLevel or AckLevel are in the future; Alert for a drifted system clock; Ask users to sync clock; Provide instructions; Exit
	// Move the AckLevel and ReadLevel after the above step to their correct timestamp => Flush to DB
	// Start jobExecutor as a go-routine
	// Start monitorAckLevelandReadLevel as a go routine
	return nil

}

func (s *scheduler) Stop() {
	// Need to understand more on the behaviour and handle accordingly
	close(s.stopChan)
}
