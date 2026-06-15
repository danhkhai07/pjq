package backqueue

import (
	"os"
	"pjq/internal/domain"
	"strconv"
	"time"
)

var (
	DEFAULT_CHAN_SIZE, _ = strconv.ParseInt(os.Getenv("DEFAULT_CHAN_SIZE"), 10, 32)
	DEFAULT_BLOCKING_TIME = 100 * time.Millisecond
)

type InProcessBackQueue struct {
	c 	chan domain.Job
}

// Blocking until receive job
func (q *InProcessBackQueue) Push(job *domain.Job) {
	q.c <- *job
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
