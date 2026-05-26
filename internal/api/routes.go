package api

import "net/http"

func addRoutes(
	svr *Server,
	mux *http.ServeMux,
) {
	mux.HandleFunc("/", svr.IndexHandler)
	mux.HandleFunc("POST /jobs", svr.PostJobHandler)
	mux.HandleFunc("GET /jobs/{id}", svr.GetJobByIDHandler)
	mux.HandleFunc("GET /jobs", svr.GetJobsByFilter)
}
