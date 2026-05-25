package queue

import (
	"testing"
	"time"

	"pjq/internal/domain"
	"pjq/internal/infra"
)

 
func TestProcessingJob(t *testing.T) {
	mockJobHandler := infra.NewMockJobHandler()
	mockStore := infra.NewMockStore(make(map[string]domain.Job))

	registry := NewRegistry()
	registry.Register("mock", mockJobHandler)

	jm := NewJobManager(NewQueue(), 3, registry, mockStore)

	go jm.Run(t.Context())

	job1 := domain.Job{
		ID: "1",
		Type: "mock",
		Status: domain.StatusPending,
		Priority: 1,
	}
	jm.PushJob(job1)

	time.Sleep(2 * time.Second)

	job, err := mockStore.Get("1")
	if err != nil {
		t.Errorf("no job pushed to store")
	}
	if job.Status != domain.StatusDone {
		t.Errorf("expected '%s', got '%s'", domain.StatusDone, job.Status)
	}
}
