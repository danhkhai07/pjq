package api

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"pjq/internal/domain"
	"pjq/internal/dto"
)

func readJSON(w http.ResponseWriter, r *http.Request, v any) bool {
	err := json.NewDecoder(r.Body).Decode(v)
	if err != nil {
		errResp := dto.ErrorResponse{
			Error: "bad request",
		}
		writeJSON(w, http.StatusBadRequest, errResp)
		return false
	}
	return true
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(v)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(`{"error":"internal server error"}`))
        return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(buf.Bytes())
}

func (svr *Server) IndexHandler(w http.ResponseWriter, r *http.Request) {
	resp := dto.IndexResponse{
		Message: "Welcome to pjq!",
	}
	writeJSON(w, http.StatusOK, resp)
}

// POST /jobs
func (svr *Server) PostJobHandler(w http.ResponseWriter, r *http.Request) {
	req := dto.PostJobRequest{}
	ok := readJSON(w, r, &req)
	if !ok {
		return
	}
	if req.IsMissingFields() {
		errResp := dto.ErrorResponse{
			Error: "bad request",
		}
		writeJSON(w, http.StatusBadRequest, errResp)
		return
	}
	
	jobID, err := svr.jobService.ProcessNewJob(req.Type, req.Payload)
	if err != nil {
		log.Print(err)
		errResp := dto.ErrorResponse{
			Error: "internal error",
		}
		writeJSON(w, http.StatusBadRequest, errResp)
		return 
	}

	resp := dto.JobIDResponse{
		ID: jobID,
	}
	writeJSON(w, http.StatusCreated, resp)
}

// GET /jobs/{id}
func (svr *Server) GetJobByIDHandler(w http.ResponseWriter, r *http.Request) {
	jobID := r.PathValue("id")
	if jobID == "" {
		errResp := dto.ErrorResponse{
			Error: "missing job id",
		}
		writeJSON(w, http.StatusBadRequest, errResp)
		return
	}

	job, err := svr.jobService.GetJobByID(jobID)
	if err != nil {
		errResp := dto.ErrorResponse{
			Error: err.Error(),
		}
		switch err {
		case domain.ErrJobNotFound:
			writeJSON(w, http.StatusNotFound, errResp)
			return
		default:
			log.Print(err)
			errResp.Error = "internal server error"
			writeJSON(w, http.StatusInternalServerError, errResp)
			return
		}
	}
	
	resp := dto.NewJobResponse(job)
	writeJSON(w, http.StatusOK, resp)
}

// GET /jobs
func (svr *Server) GetJobsByFilter(w http.ResponseWriter, r *http.Request) {
	filter := domain.JobFilter{}
	if r.Header.Get("Content-Type") == "application/json" {
		ok := readJSON(w, r, &filter)
		if !ok {
			return
		}
	}

	jobs, err := svr.jobService.ListJobsWithFilter(filter)
	if err != nil {
		errResp := dto.ErrorResponse{
			Error: err.Error(),
		}
		switch err {
		default:
			log.Print(err)
			errResp.Error = "internal server error"
			writeJSON(w, http.StatusInternalServerError, errResp)
			return
		}
	}

	jobsResponse := make([]dto.JobResponse, len(jobs))
	for i, job := range jobs {
		jobsResponse[i] = dto.NewJobResponse(job)
	}

	resp := dto.ListJobsResponse{
		Jobs: jobsResponse,
		Total: len(jobsResponse),
	}
	writeJSON(w, http.StatusOK, resp)
}
