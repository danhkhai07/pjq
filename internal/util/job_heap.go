package util

import (
	"pjq/internal/domain"
)

type JobHeap []domain.Job

func (h JobHeap) Len() int 			{ return len(h) }
func (h JobHeap) Swap(i, j int) 	{ h[i], h[j] = h[j], h[i] }
func (h *JobHeap) Push(x any) 		{ *h = append(*h, x.(domain.Job))}

func (h JobHeap) Peek() *domain.Job {
	if h.Len() < 1 {
		return nil
	}
	return &h[0]
}

func (h JobHeap) Less(i, j int) bool {
	if *h[i].RunAt != *h[j].RunAt {
		return h[i].RunAt.Before(*h[j].RunAt)
	}
	if h[i].Priority != h[j].Priority {
		return h[i].Priority > h[j].Priority
	}
	return h[i].CreatedAt.Before(h[j].CreatedAt)
}

func (h *JobHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}
