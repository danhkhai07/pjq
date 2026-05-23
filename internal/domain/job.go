package domain

import "time"

type Status string

const (
	Pending 	Status = "pending"
	Running		Status = "running"
	Done		Status = "done"
	Failed		Status = "failed"
	Retrying	Status = "retrying"
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
