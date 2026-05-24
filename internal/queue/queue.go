package queue

import (
	"sync"
)

type Queue[T any] struct {
	mu  sync.Mutex
	arr []T
	priority []int
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{
		arr: make([]T, 0),
	}
}

func (q *Queue[T]) Push(v T) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.arr = append(q.arr, v)
}

func (q *Queue[T]) Pop() (T, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.arr) == 0 {
		var zero T
		return zero, false
	}
	head := q.arr[0]
	q.arr = q.arr[1:]
	return head, true
}

func (q *Queue[T]) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.arr)
}

func (q *Queue[T]) Snapshot() []T {
	q.mu.Lock()
	defer q.mu.Unlock()
	cp := make([]T, len(q.arr))
	copy(cp, q.arr)
	return cp
}
