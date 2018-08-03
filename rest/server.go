package rest

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/almostmoore/gbquestion/storage"
	"github.com/almostmoore/gbquestion/vars"
)

// Server struct represens a server
type Server struct {
	addr string
	qs   *storage.QuestionStorage
}

// NewServer creates a new server
func NewServer(addr string, qs *storage.QuestionStorage) *Server {
	return &Server{
		addr: addr,
		qs:   qs,
	}
}

// Run starts a server
func (s *Server) Run() {
	r := mux.NewRouter()

	r.HandleFunc("/version", s.versionHandler).Methods(http.MethodGet)
	r.HandleFunc("/", s.insertQuestionHandler).Methods(http.MethodPost)
	r.HandleFunc("/{id:[0-9]+}", s.getQuestionHandler).Methods(http.MethodGet)
	r.HandleFunc("/{id:[0-9]+}", s.updateQuestionHandler).Methods(http.MethodPut)
	r.HandleFunc("/{id:[0-9]+}", s.deleteQuestionHandler).Methods(http.MethodDelete)

	loggedRouter := handlers.LoggingHandler(os.Stdout, r)

	http.ListenAndServe(s.addr, loggedRouter)
}

func (s *Server) versionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	encoder := json.NewEncoder(w)
	encoder.Encode(map[string]string{
		"version": vars.Version,
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

	var q storage.Question
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

func (s *Server) updateQuestionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	decoder := json.NewDecoder(r.Body)
	vars := mux.Vars(r)

	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(map[string]string{
			"error": err.Error(),
		})

		return
	}

	var q storage.Question
	err = decoder.Decode(&q)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(map[string]string{
			"error": err.Error(),
		})

		return
	}

	q.ID = id
	_, err = s.qs.Put(q)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(map[string]string{
			"error": err.Error(),
		})

		return
	}

	w.WriteHeader(http.StatusOK)
	encoder.Encode(q)
}

func (s *Server) deleteQuestionHandler(w http.ResponseWriter, r *http.Request) {
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

	err = s.qs.Delete(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(map[string]string{
			"error": err.Error(),
		})

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
