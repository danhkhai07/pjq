package queue

import (
	"sync"
	"container/heap"

	"pjq/internal/domain"
	"pjq/internal/util"
)

type FrontQueue struct {
	mu  sync.Mutex
	heap *util.JobHeap
}

func NewQueue() *FrontQueue {
	return &FrontQueue{
		heap: &util.JobHeap{},
	}
}

func (q *FrontQueue) Push(job domain.Job) {
	q.mu.Lock()
	defer q.mu.Unlock()
	heap.Push(q.heap, job)
}

func (q *FrontQueue) Pop() (domain.Job, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.heap.Len() == 0 {
		var zero domain.Job
		return zero, false
	}
	head := heap.Pop(q.heap).(domain.Job)
	return head, true
}

func (q *FrontQueue) Peek() *domain.Job {
	return q.heap.Peek()
}

func (q *FrontQueue) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.heap.Len()
}

func (q *FrontQueue) Snapshot() []domain.Job {
	q.mu.Lock()
	defer q.mu.Unlock()
	cp := make([]domain.Job, len(*q.heap))
	copy(cp, *q.heap)
	return cp
}
