package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

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

	r.HandleFunc("/", s.filterQuestionHandler).Methods(http.MethodGet)
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

	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		s.sendError(w, http.StatusInternalServerError, err)
		return
	}

	q, err := s.qs.Get(id)
	if err != nil {
		s.sendError(w, http.StatusInternalServerError, err)
		return
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(q)
}

func (s *Server) insertQuestionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	var q storage.Question
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&q)
	if err != nil {
		s.sendError(w, http.StatusInternalServerError, err)
		return
	}

	q.ID = 0
	id, err := s.qs.Put(q)
	if err != nil {
		s.sendError(w, http.StatusInternalServerError, err)
		return
	}

	q.ID = id

	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	encoder.Encode(q)
}

func (s *Server) updateQuestionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	vars := mux.Vars(r)

	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		s.sendError(w, http.StatusInternalServerError, err)
		return
	}

	var q storage.Question
	err = decoder.Decode(&q)
	if err != nil {
		s.sendError(w, http.StatusInternalServerError, err)
		return
	}

	q.ID = id
	_, err = s.qs.Put(q)
	if err != nil {
		s.sendError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	encoder.Encode(q)
}

func (s *Server) deleteQuestionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		s.sendError(w, http.StatusInternalServerError, err)
		return
	}

	err = s.qs.Delete(id)
	if err != nil {
		s.sendError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) filterQuestionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	fmt.Println(r.URL.Query().Get("ignore"))

	limit, _ := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 32)
	active := r.URL.Query().Get("active")

	filter := &storage.QuestionFilter{
		Limit:    int(limit),
		IsActive: active == "1" || active == "true",
	}

	ignoreStr := strings.Split(r.URL.Query().Get("ignore"), ",")
	for i := 0; i < len(ignoreStr); i++ {
		id, _ := strconv.ParseUint(ignoreStr[i], 10, 64)
		if id != 0 {
			filter.IgnoreIds = append(filter.IgnoreIds, id)
		}
	}

	fmt.Printf("%+v", filter)
	questions, err := s.qs.Filter(filter)
	if err != nil {
		s.sendError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	encoder.Encode(questions)
}

func (s *Server) sendError(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)
	encoder := json.NewEncoder(w)
	encoder.Encode(map[string]string{
		"error": err.Error(),
	})
}
