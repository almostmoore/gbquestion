package main

import (
	"encoding/json"

	"github.com/boltdb/bolt"
)

var questionsBucketName = []byte("questions")

// Question structure
type Question struct {
	ID       uint64 `json:"id"`
	Text     string `json:"text"`
	IsGood   bool   `json:"is_good"`
	IsActive bool   `json:"is_active"`
}

// QuestionStorage stores questions
type QuestionStorage struct {
	db *bolt.DB
}

// NewQuestionStorage creates a new question storage
func NewQuestionStorage(DB *bolt.DB) *QuestionStorage {
	return &QuestionStorage{
		db: DB,
	}
}

// Insert creates a new question into db
func (qs *QuestionStorage) Insert(q Question) (uint64, error) {
	var id uint64

	err := qs.db.Batch(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(questionsBucketName)
		if err != nil {
			return err
		}

		id, err = b.NextSequence()
		if err != nil {
			return err
		}

		q.ID = id
		data, err := json.Marshal(q)
		if err != nil {
			return err
		}

		return b.Put(uinttob(id), data)
	})

	return id, err
}

// Update renews a question in the database
func (qs *QuestionStorage) Update(q Question) error {
	data, err := json.Marshal(q)
	if err != nil {
		return err
	}

	id := uinttob(q.ID)

	return qs.db.Batch(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(questionsBucketName)
		if err != nil {
			return err
		}

		return b.Put(id, data)
	})
}

// Get returns a question by it's ID
func (qs *QuestionStorage) Get(id uint64) (Question, error) {
	var q Question

	err := qs.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(questionsBucketName)

		data := b.Get(uinttob(id))
		return json.Unmarshal(data, &q)
	})

	return q, err
}
