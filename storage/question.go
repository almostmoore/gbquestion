package storage

import (
	"encoding/binary"
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

// QuestionFilter struct used for question filtration
type QuestionFilter struct {
	IsActive  bool
	Limit     int
	IgnoreIds []uint64
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

// Delete func removes question by id
func (qs *QuestionStorage) Delete(id uint64) error {
	return qs.db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket(questionsBucketName)
		return b.Delete(utils.Uinttob(id))
	})
}

// Filter func searches questions by filter
func (qs *QuestionStorage) Filter(filter *QuestionFilter) ([]Question, error) {
	questions := make([]Question, 0, filter.Limit)

	ignoreIds := make(map[uint64]bool, len(filter.IgnoreIds))
	for i := 0; i < len(filter.IgnoreIds); i++ {
		ignoreIds[filter.IgnoreIds[i]] = true
	}

	err := qs.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(questionsBucketName)
		c := b.Cursor()

		for k, v := c.First(); k != nil && len(questions) < filter.Limit; k, v = c.Next() {
			uintKey := binary.BigEndian.Uint64(k)
			if ignoreIds[uintKey] {
				continue
			}

			var q Question
			err := json.Unmarshal(v, &q)
			if err != nil {
				return err
			}

			if q.IsActive != filter.IsActive {
				continue
			}

			questions = append(questions, q)
		}

		return nil
	})

	return questions, err
}
