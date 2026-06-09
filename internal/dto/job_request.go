package dto

import (
	"time"
	"encoding/json"

	"pjq/internal/domain"
)

type PostJobRequest struct {
	Type 			string 			`json:"type"`
	Payload 		json.RawMessage	`json:"payload"`
	RunAt			*time.Time		`json:"run_at,omitempty"`
}

func (req *PostJobRequest) IsMissingFields() bool {
	if req.Type == "" || req.Payload == nil {
		return true
	}
	return false
}

type GetJobsWithFilterRequest struct {
	Status			domain.Status	`json:"status,omitempty"`
	Type 			string			`json:"type,omitempty"`
	Retriable 		bool			`json:"retriable,omitempty"`
}

func (req * GetJobsWithFilterRequest) IsMissingFields() bool {
	return true
}
