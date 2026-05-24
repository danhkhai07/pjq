package queue

import (
	"sync"
	"container/heap"

	"pjq/internal/domain"
	"pjq/internal/util"
)

type Queue struct {
	mu  sync.Mutex
	heap *util.JobHeap
}

func NewQueue() *Queue {
	return &Queue{
		heap: &util.JobHeap{},
	}
}

func (q *Queue) Push(job domain.Job) {
	q.mu.Lock()
	defer q.mu.Unlock()
	heap.Push(q.heap, job)
}

func (q *Queue) Pop() (domain.Job, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.heap.Len() == 0 {
		var zero domain.Job
		return zero, false
	}
	head := heap.Pop(q.heap).(domain.Job)
	return head, true
}

func (q *Queue) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.heap.Len()
}

func (q *Queue) Snapshot() []domain.Job {
	q.mu.Lock()
	defer q.mu.Unlock()
	cp := make([]domain.Job, len(*q.heap))
	copy(cp, *q.heap)
	return cp
}
