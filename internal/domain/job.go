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
	Result any
	Priority int
	Retries int
	MaxRetries int
	CreatedAt time.Time
	StartedAt *time.Time
	FinishedAt *time.Time
	RunAt *time.Time
	Error string
	Logs []string
}

func NewJob(
	id string,
	jobType string,
	payload []byte,
	runAt *time.Time,
	priority int,
	maxRetries int,
) Job {
	now := time.Now()
	job := Job{
		ID: id,
		Type: jobType,
		Payload: payload,
		Status: StatusPending,
		Result: nil,
		Priority: priority,
		Retries: 0,
		MaxRetries: maxRetries,
		CreatedAt: now,
		StartedAt: nil,
		FinishedAt: nil,
		RunAt: runAt,
		Error: "",
		Logs: make([]string, 0, 10),
	}
	if job.RunAt == nil {
		job.RunAt = &now
	}
	return job
}
