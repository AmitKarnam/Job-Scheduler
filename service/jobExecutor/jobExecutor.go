package jobexecutor

import (
	"fmt"

	"github.com/AmitKarnam/Job-Scheduler/models"
)

// JobExecutor defines the interface for executing a job. Concrete implementations will handle the specifics of executing different job types. The Execute method should perform the job's action and return an error if execution fails.
type JobExecutor interface {
	Execute(job *models.Job) error
}

// executionRegistry maps JobType to their corresponding JobExecutor implementations. This allows for dynamic execution of different job types based on their defined executors. ( Strategy Pattern )
var ExecutionRegistry = map[models.JobType]JobExecutor{}

// AddToRegistry allows for registering a JobExecutor for a specific JobType. This function can be called during application initialization to set up the necessary executors for the job types that will be used in the system.
func AddToRegistry(jobType models.JobType, jobExecutor JobExecutor) {
	if _, exists := ExecutionRegistry[jobType]; exists {
		fmt.Printf("executor for job type %s already exists", jobType)
		return
	}
	ExecutionRegistry[jobType] = jobExecutor
}

type EmailJobExecutor struct{}

func (e EmailJobExecutor) Execute(job *models.Job) error {
	// extract payload
	// send email
	return nil
}

type MobileNotificationJobExecutor struct{}

func (m MobileNotificationJobExecutor) Execute(job *models.Job) error {
	// extract payload
	// send mobile notification
	return nil
}
