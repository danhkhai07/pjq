package infra 

import (
	"time"
	"pjq/internal/domain"
)

var (
	DEFAULT_BLOCKING_TIME = 100 * time.Millisecond
)

type InProcessBackQueue struct {
	c 	chan domain.Job
}

func NewInProcessBackQueue(size int) *InProcessBackQueue {
	return &InProcessBackQueue{
		c: make(chan domain.Job, size),
	}
}

// Blocking until receive job
func (q *InProcessBackQueue) Push(job *domain.Job) error {
	q.c <- *job
	return nil
}

// Wait for 100ms before returning false
func (q *InProcessBackQueue) Pop() (*domain.Job, bool) {
	timer := time.NewTimer(DEFAULT_BLOCKING_TIME)
	select {
	case job := <-q.c:
		return &job, true
	case <-timer.C:
		return nil, false
	}
}

func (q *InProcessBackQueue) GetName() string { return "In-process Back Queue"}
