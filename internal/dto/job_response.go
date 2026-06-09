package dto

import (
	"time"

	"pjq/internal/domain"
)

type JobIDResponse struct {
	ID 			string			`json:"id"`
}

type JobResponse struct {
	ID 			string			`json:"id"`
	Type 		string			`json:"type"`
	Payload 	[]byte			`json:"payload"`
	Result		any				`json:"result,omitempty"`
	Status 		domain.Status	`json:"status"`	
	Priority 	int				`json:"priority"`
	Retries 	int				`json:"retries"`
	MaxRetries 	int				`json:"max_retries"`
	CreatedAt 	time.Time		`json:"created_at"`
	StartedAt 	*time.Time		`json:"started_at,omitempty"`
	FinishedAt 	*time.Time		`json:"finished_at,omitempty"`
	RunAt		*time.Time		`json:"run_at,omitempty"`
	Error 		string			`json:"error,omitempty"`
	Logs 		[]string		`json:"logs,omitempty"`
}

func NewJobResponse(job domain.Job) JobResponse {
	return JobResponse{
		ID: 			job.ID,
		Type: 			job.Type,
		Payload: 		job.Payload,
		Status: 		job.Status,
		Result:			job.Result,
		Priority: 		job.Priority,
		Retries: 		job.Retries,
		MaxRetries: 	job.MaxRetries,
		CreatedAt: 		job.CreatedAt,
		StartedAt: 		job.StartedAt,
		FinishedAt: 	job.FinishedAt,
		RunAt: 			job.RunAt,		
		Error: 			job.Error,
		Logs: 			job.Logs,
	}
}


type ListJobsResponse struct {
	Jobs 		[]JobResponse	`json:"jobs"`
	Total		int 			`json:"total"`
}
