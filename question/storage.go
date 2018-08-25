package question

import (
	"encoding/binary"

	"github.com/almostmoore/gbquestion/utils"
	"github.com/boltdb/bolt"
	"github.com/golang/protobuf/proto"
)

var questionsBucketName = []byte("questions")

// Storage stores questions
type Storage struct {
	db *bolt.DB
}

// NewStorage creates a new question storage
func NewStorage(DB *bolt.DB) *Storage {
	return &Storage{
		db: DB,
	}
}

// Put creates or updates a question into db
func (qs *Storage) Put(q Question) (uint64, error) {
	id := q.Id

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

		q.Id = id
		data, err := proto.Marshal(&q)
		if err != nil {
			return err
		}

		return b.Put(utils.Uinttob(id), data)
	})

	return id, err
}

// Get returns a question by it's ID
func (qs *Storage) Get(id uint64) (*Question, error) {
	q := &Question{}

	err := qs.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(questionsBucketName)

		data := b.Get(utils.Uinttob(id))
		return proto.Unmarshal(data, q)
	})

	return q, err
}

// Delete func removes question by id
func (qs *Storage) Delete(id uint64) error {
	return qs.db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket(questionsBucketName)
		return b.Delete(utils.Uinttob(id))
	})
}

// Filter func searches questions by filter
func (qs *Storage) Filter(filter *Filter) ([]*Question, error) {
	questions := make([]*Question, 0, filter.Limit)

	ignoreIds := make(map[uint64]bool, len(filter.IgnoreIds))
	for i := 0; i < len(filter.IgnoreIds); i++ {
		ignoreIds[filter.IgnoreIds[i]] = true
	}

	err := qs.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(questionsBucketName)
		c := b.Cursor()
		var offset int32

		for k, v := c.First(); k != nil && int32(len(questions)) < filter.Limit; k, v = c.Next() {
			uintKey := binary.BigEndian.Uint64(k)
			if ignoreIds[uintKey] {
				continue
			}

			q := &Question{}
			err := proto.Unmarshal(v, q)
			if err != nil {
				return err
			}

			if q.IsActive != filter.IsActive {
				continue
			}

			if offset < filter.Offset {
				offset++
				continue
			}

			questions = append(questions, q)
		}

		return nil
	})

	return questions, err
}
