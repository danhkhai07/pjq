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
	Status 		domain.Status	`json:"status"`	
	Priority 	int				`json:"priority"`
	Retries 	int				`json:"retries"`
	MaxRetries 	int				`json:"max_retries"`
	CreatedAt 	time.Time		`json:"created_at"`
	StartedAt 	time.Time		`json:"started_at"`
	FinishedAt 	time.Time		`json:"finished_at"`
	Error 		string			`json:"error"`
	Logs 		[]string		`json:"logs"`
}


type ListJobsResponse struct {
	Jobs 		[]JobResponse	`json:"jobs"`
	Total		int 			`json:"total"`
}
