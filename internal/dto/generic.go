package dto

import (
	
)

type StatusResponse struct {
	Status 			string `json:"status"`
	Message			string `json:"message"`
}

type ErrorResponse struct {
	Error 			string `json:"status"`
	Message			string `json:"message"`
}
