package infra

import (
	"context"
	"errors"
	"fmt"
	"time"

	"pjq/internal/domain"
)

type MockJobHandler struct {}

func NewMockJobHandler() *MockJobHandler {
	return &MockJobHandler{}
}

func (mjh *MockJobHandler) Handle(ctx context.Context, job *domain.Job, log func(string)) error {
	fmt.Printf("Handler: handling job %v\n", job.ID)
	time.Sleep(1 * time.Second)
	fmt.Printf("Handler: job failed %v\n", job.ID)
	return errors.New("mock error")
}
