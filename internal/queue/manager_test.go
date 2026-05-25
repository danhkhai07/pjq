package queue

import (
	"context"
	"testing"
	"time"

	"pjq/internal/domain"
	"pjq/internal/infra"
)

func testingEnv() (*JobManager, domain.JobStore, context.CancelFunc) {
	mockJobHandler := infra.NewMockJobHandler()
	mockStore := infra.NewMockStore(make(map[string]domain.Job))

	registry := NewRegistry()
	registry.Register("mock", mockJobHandler)

	jm := NewJobManager(NewQueue(), 3, registry, mockStore)

	ctx, cancel := context.WithCancel(context.Background())
	go jm.Run(ctx)
	return jm, mockStore, cancel
}

func TestProcessingJob(t *testing.T) {
	jm, store, ctxCancel := testingEnv()
	defer ctxCancel()

	go jm.Run(t.Context())

	job1 := domain.NewJob("1", "mock", []byte{}, 1, 0)	
	jm.PushJob(*job1)

	time.Sleep(2 * time.Second)

	job, err := store.Get("1")
	if err != nil {
		t.Errorf("no job pushed to store")
	}
	if job.Status != domain.StatusFailed {
		t.Errorf("expected '%s', got '%s'", domain.StatusFailed, job.Status)
	}
}

func TestRetryAlgorithm(t *testing.T) {
	jm, store, ctxCancel := testingEnv()
	defer ctxCancel()
	
	job1 := domain.NewJob("1", "mock", []byte{}, 1, 3)	
	jm.PushJob(*job1)

	time.Sleep(5 * time.Second)
	job, err := store.Get("1")
	if err != nil {
		t.Errorf("no job pushed to store")
	}
	if job.Status != domain.StatusFailed {
		t.Errorf("Status: expected '%s', got '%s'", domain.StatusFailed, job.Status)
	}
	if job.Retries != 3 {
		t.Errorf("Status: expected '%v', got '%v'", 3, job.Retries)
	}
}
