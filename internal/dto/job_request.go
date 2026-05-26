package dto

import (
	"encoding/json"

	"pjq/internal/domain"
)

type PostJobRequest struct {
	Type 			string 			`json:"type"`
	Payload 		json.RawMessage	`json:"payload"`
}

func (req *PostJobRequest) IsMissingFields() bool {
	if req.Type == "" || req.Payload == nil {
		return false
	}
	return true
}

type GetJobsWithFilterRequest struct {
	Status			domain.Status	`json:"status,omitempty"`
	Type 			string			`json:"type,omitempty"`
	Retriable 		bool			`json:"retriable,omitempty"`
}

func (req * GetJobsWithFilterRequest) IsMissingFields() bool {
	return true
}
