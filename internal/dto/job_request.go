package dto

import (
	"encoding/json"

	"pjq/internal/domain"
)

type PostJobRequest struct {
	Type 			string 			`json:"type"`
	Payload 		json.RawMessage	`json:"payload"`
}

type GetJobsWithFilterRequest struct {
	Status			domain.Status	`json:"status,omitempty"`
	Type 			string			`json:"type,omitempty"`
	Retriable 		bool			`json:"retriable,omitempty"`
}
