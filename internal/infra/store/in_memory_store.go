package infra

import (
	"errors"
	"pjq/internal/domain"
)

type InMemoryStore struct {
	data map[string]domain.Job
}

func NewInMemoryStore(data map[string]domain.Job) *InMemoryStore {
	return &InMemoryStore{
		data: data,
	}
}

func (store *InMemoryStore) Save(job domain.Job) error {
	store.data[job.ID] = job
	return nil
}

func (store *InMemoryStore) Get(id string) (domain.Job, error) {
	job, ok := store.data[id]
	if !ok {
		return job, errors.New("key not found")
	}
	return job, nil
}
 
func (store *InMemoryStore) Delete(id string) error {
	delete(store.data, id)
	return nil
}
 
// List accordingly to the filter. O(n) operation.
func (store *InMemoryStore) List(filter domain.JobFilter) ([]domain.Job, error) {
	var result []domain.Job
	for _, job := range store.data {
		if filter.Pass(job) {
			result = append(result, job)
		}
	}
	return result, nil
}
