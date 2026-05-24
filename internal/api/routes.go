package api

import "net/http"

func addRoutes(
	svr *Server,
	mux *http.ServeMux,
) {
	mux.HandleFunc("/", svr.indexHandler)
}
