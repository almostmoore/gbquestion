package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	addr string
	qs   *QuestionStorage
}

func NewServer(addr string, qs *QuestionStorage) *Server {
	return &Server{
		addr: addr,
		qs:   qs,
	}
}

func (s *Server) Run() {
	r := mux.NewRouter()

	r.HandleFunc("/version", versionHandler).Methods(http.MethodGet)

	http.ListenAndServe(s.addr, r)
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	encoder := json.NewEncoder(w)
	encoder.Encode(map[string]string{
		"version": version,
	})
}
