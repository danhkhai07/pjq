package queue

import (
	"context"
	"fmt"
	"testing"
	"time"

	"pjq/internal/domain"
)

type MockJobHandler struct {}

func (mjh *MockJobHandler) Handle(ctx context.Context, job *domain.Job, log func(string)) error {
	fmt.Printf("Handler: handling job %v\n", job.ID)
	time.Sleep(1 * time.Second)
	job.Status = domain.Done
	fmt.Printf("Handler: done job %v\n", job.ID)
	return nil
}
 
func TestProcessingJob(t *testing.T) {
	mockRegistry := NewRegistry()
	var mockJobHandler MockJobHandler
	mockRegistry.Register("mock", &mockJobHandler)

	jm := NewJobManager(NewQueue(), 3, mockRegistry)

	ctx, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()
	go jm.Run(ctx)

	job1 := domain.Job{
		ID: "1",
		Type: "mock",
		Status: domain.Pending,
		Priority: 1,
	}
	jm.PushJob(job1)

	time.Sleep(2 * time.Second)
}
