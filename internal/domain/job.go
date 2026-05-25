package domain

import "time"

type Status string

const (
	StatusPending 		Status = "pending"
	StatusRunning		Status = "running"
	StatusDone			Status = "done"
	StatusFailed		Status = "failed"
	StatusRetrying		Status = "retrying"
)

type Job struct {
	ID string
	Type string
	Payload []byte
	Status Status
	Priority int
	Retries int
	MaxRetries int
	CreatedAt time.Time
	StartedAt time.Time
	FinishedAt time.Time
	Error string
	Logs []string
}

func NewJob(
	id string,
	jobType string,
	payload []byte,
	priority int,
	maxRetries int,
) *Job {
	job := Job{
		ID: id,
		Type: jobType,
		Payload: payload,
		Status: StatusPending,
		Priority: priority,
		Retries: 0,
		MaxRetries: maxRetries,
		CreatedAt: time.Now(),
		StartedAt: time.Time{},
		FinishedAt: time.Time{},
		Error: "",
		Logs: nil,
	}
	return &job
}
