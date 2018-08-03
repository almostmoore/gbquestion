package storage

import (
	"encoding/json"

	"github.com/almostmoore/gbquestion/utils"
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

// Put creates or updates a question into db
func (qs *QuestionStorage) Put(q Question) (uint64, error) {
	id := q.ID

	err := qs.db.Batch(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(questionsBucketName)
		if err != nil {
			return err
		}

		if id == 0 {
			id, err = b.NextSequence()
			if err != nil {
				return err
			}
		}

		q.ID = id
		data, err := json.Marshal(q)
		if err != nil {
			return err
		}

		return b.Put(utils.Uinttob(id), data)
	})

	return id, err
}

// Get returns a question by it's ID
func (qs *QuestionStorage) Get(id uint64) (Question, error) {
	var q Question

	err := qs.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(questionsBucketName)

		data := b.Get(utils.Uinttob(id))
		return json.Unmarshal(data, &q)
	})

	return q, err
}
