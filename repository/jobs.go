package repository

import (
	"fmt"
	"sync"

	"github.com/AmitKarnam/Job-Scheduler/models"
)

// JobStore defines the contract for job storage
type JobRepo interface {
	Save(job *models.Job) error
	Get(jobID string) (*models.Job, error)
	GetAll() ([]*models.Job, error)
	Update(job *models.Job) error
	Delete(jobID string) error
	GetByStatus(status models.Status) ([]*models.Job, error)
}

// InMemoryJobStore is a thread-safe in-memory job store
type InMemoryJobStore struct {
	mu   sync.RWMutex
	jobs map[string]*models.Job
}

func NewInMemoryJobStore() JobRepo {
	return &InMemoryJobStore{
		jobs: make(map[string]*models.Job),
	}
}

func (js *InMemoryJobStore) Save(job *models.Job) error {
	js.mu.Lock()
	defer js.mu.Unlock()

	if _, exists := js.jobs[job.ID]; exists {
		return fmt.Errorf("job %s already exists", job.ID)
	}
	js.jobs[job.ID] = job
	return nil
}

func (js *InMemoryJobStore) Get(jobID string) (*models.Job, error) {
	js.mu.RLock()
	defer js.mu.RUnlock()

	job, exists := js.jobs[jobID]
	if !exists {
		return nil, fmt.Errorf("job %s not found", jobID)
	}
	return job, nil
}

func (js *InMemoryJobStore) GetAll() ([]*models.Job, error) {
	js.mu.RLock()
	defer js.mu.RUnlock()

	jobs := make([]*models.Job, 0, len(js.jobs))
	for _, job := range js.jobs {
		jobs = append(jobs, job)
	}
	return jobs, nil
}

func (js *InMemoryJobStore) Update(job *models.Job) error {
	js.mu.Lock()
	defer js.mu.Unlock()

	if _, exists := js.jobs[job.ID]; !exists {
		return fmt.Errorf("job %s not found", job.ID)
	}
	js.jobs[job.ID] = job
	return nil
}

func (js *InMemoryJobStore) Delete(jobID string) error {
	js.mu.Lock()
	defer js.mu.Unlock()

	delete(js.jobs, jobID)
	return nil
}

func (js *InMemoryJobStore) GetByStatus(status models.Status) ([]*models.Job, error) {
	js.mu.RLock()
	defer js.mu.RUnlock()

	jobs := make([]*models.Job, 0)
	for _, job := range js.jobs {
		if job.Status == status {
			jobs = append(jobs, job)
		}
	}
	return jobs, nil
}
