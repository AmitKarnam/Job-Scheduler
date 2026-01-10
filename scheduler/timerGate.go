package scheduler

import (
	"fmt"
	"sync"
	"time"
)

type JobType string

const (
	EmailJob           JobType = "email_job"
	MobileNotification JobType = "mobile_notification"
)

type Job struct {
	jobType       JobType
	executionTime time.Time
}

func (j *Job) Execute() {
	fmt.Println("Executing Job: ", j.jobType)
}

func main() {

	// WaitGroup for waiting on the timer gate go routine.
	var wg sync.WaitGroup

	// Dummy Heap
	min_heap := make([]Job, 1)

	// Signal Interrupt for new job
	newJobChan := make(chan Job)

	// Stop chan to execution of the timer gate
	stopChan := make(chan bool)

	// Initialise Timer
	timer := time.NewTimer(0)
	if !timer.Stop() {
		select {
		case <-timer.C:
		default:
		}
	}

	// Dummy Time initialised and put into the heap
	jobExecTime, err := time.Parse(time.RFC3339, "2025-12-17T18:10:00Z")
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	// Dummy Job Definition
	dummyJob := Job{
		jobType:       EmailJob,
		executionTime: jobExecTime,
	}

	min_heap = append(min_heap, dummyJob)

	// Logic to calculate the next job and set timer

	// Handling Timer Gate

	// The new job should only be sent into the newJob channel if the job is within the execution ExecTime > AckLevel AND ExecTime < ReadLevelTime

	// Timer gate logic
	wg.Add(1)
	go func() {
		select {
		case <-timer.C:
			jobToBeExecuted := min_heap[0]
			// Async Execution: Use Worker Pool method for Jobs; Can use in memory channel for now, which can be replaced with Message Queue later for Durable Exection.
			go jobToBeExecuted.Execute()

			// Update Heap
			// Logic to update the heap? OR Remove the executed job and update the heap?

			// Calculate diff duration
			// diffDuration := time.Now().UTC() - nextJob_NextExecutionTime

			// Set timer for the next job
			// if !timer.Reset(diffDuration) {
			// 	fmt.Print("Error updating timer")
			// 	stopChan <- true
			// }

		case newJob := <-newJobChan:
			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}

			// Insert the job into heap
			min_heap = append(min_heap, newJob)

		case <-stopChan:
			timer.Stop()
			wg.Done()
			return

		}
	}()

	wg.Wait()

}
