package question

import (
	fmt "fmt"

	context "golang.org/x/net/context"
)

// RPCService is a simple grpc question service
type RPCService struct {
	storage *Storage
}

// NewRPCService returns a new service
func NewRPCService(s *Storage) *RPCService {
	return &RPCService{
		storage: s,
	}
}

// List func returns a filtered list of questions
func (s RPCService) List(ctx context.Context, filter *Filter) (*QuestionList, error) {
	list, err := s.storage.Filter(filter)
	if err != nil {
		return nil, fmt.Errorf("Coudln't get questions from the storage: %v", err)
	}

	return &QuestionList{
		Questions: list,
	}, nil
}

// Put func saves a question
func (s RPCService) Put(ctx context.Context, q *Question) (*Question, error) {
	id, err := s.storage.Put(*q)
	if err != nil {
		return nil, fmt.Errorf("Couldn't save a message: %v", err)
	}

	q.Id = id
	return q, nil
}

// Get func returns a question by ID
func (s RPCService) Get(ctx context.Context, req *IdRequest) (*Question, error) {
	return s.storage.Get(req.Id)
}

// Delete func delete question by ID
func (s RPCService) Delete(ctx context.Context, req *IdRequest) (*Void, error) {
	return &Void{}, s.storage.Delete(req.Id)
}
