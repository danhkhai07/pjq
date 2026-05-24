package api

import (
	"encoding/json"
	"net/http"

	"pjq/internal/dto"
)

func (svr *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	resp := dto.IndexResponse{
		Message: "Welcome to pjq!",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
