package infra

import (
	"errors"
	"pjq/internal/domain"
)

type MockStore struct {
	mp map[string]domain.Job
}

func NewMockStore(data map[string]domain.Job) *MockStore {
	return &MockStore{
		mp: data,
	}
}

func (ms *MockStore) Save(job domain.Job) error {
	ms.mp[job.ID] = job
	return nil
}

func (ms *MockStore) Get(id string) (domain.Job, error) {
	job, ok := ms.mp[id]
	if !ok {
		return job, errors.New("key not found")
	}
	return job, nil
}
 
func (ms *MockStore) Delete(id string) error {
	delete(ms.mp, id)
	return nil
}
 
// List accordingly to the filter. O(n) operation.
func (ms *MockStore) List(filter domain.JobFilter) ([]domain.Job, error) {
	var result []domain.Job
	for _, job := range ms.mp {
		if filter.Pass(job) {
			result = append(result, job)
		}
	}
	return result, nil
}
