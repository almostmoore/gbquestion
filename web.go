package main

import (
	"encoding/json"
	"net/http"
	"strconv"

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

	r.HandleFunc("/version", s.versionHandler).Methods(http.MethodGet)
	r.HandleFunc("/", s.insertQuestionHandler).Methods(http.MethodPost)
	r.HandleFunc("/{id:[0-9]+}", s.getQuestionHandler).Methods(http.MethodGet)

	http.ListenAndServe(s.addr, r)
}

func (s *Server) versionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	encoder := json.NewEncoder(w)
	encoder.Encode(map[string]string{
		"version": version,
	})
}

func (s *Server) getQuestionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(w)

	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	q, err := s.qs.Get(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	encoder.Encode(q)
}

func (s *Server) insertQuestionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(w)

	var q Question
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&q)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	q.ID = 0
	id, err := s.qs.Put(q)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	q.ID = id

	w.WriteHeader(http.StatusOK)
	encoder.Encode(q)
}
