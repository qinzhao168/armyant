package storage

/*import (
	"fmt"

	"github.com/yiqinguo/armyant/pkg/models"

	"github.com/coreos/bbolt"
	"github.com/golang/protobuf/proto"
)

var _ Storage = &bboltStorage{}

const (
	jobBucket string = "armyant"
)

type bboltStorage struct {
	//dir string
	db *bolt.DB
}

func NewBboltStorage(dir string) (*bboltStorage, error) {
	db, err := bolt.Open(dir, 0666, nil)
	if err != nil {
		return nil, err
	}
	return &bboltStorage{
		db: db,
	}, err
}

func (b *bboltStorage) Add(job models.Job) error {
	body, err := proto.Marshal(&job)
	if err != nil {
		return err
	}
	tx, err := b.db.Begin(true)
	if err != nil {
		return err
	}
	bucket, err := tx.CreateBucketIfNotExists([]byte(jobBucket))
	if err != nil {
		tx.Rollback()
		return err
	}
	err = bucket.Put([]byte(job.Id), body)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (b *bboltStorage) Get(id string) (models.Job, error) {
	var job models.Job
	tx, err := b.db.Begin(false)
	if err != nil {
		return job, err
	}
	bucket := tx.Bucket([]byte(jobBucket))
	body := bucket.Get([]byte(id))

	err = proto.UnmarshalMerge(body, &job)
	if err != nil {
		return job, err
	}
	return job, nil
}

func (b *bboltStorage) List() ([]models.Job, error) {
	var jobs []models.Job
	tx, err := b.db.Begin(false)
	if err != nil {
		return jobs, err
	}
	bucket := tx.Bucket([]byte(jobBucket))
	err = bucket.ForEach(func(k, v []byte) error {
		fmt.Printf("A %s is %s.\n", k, v)
		return nil
	})
	return []models.Job{}, nil
}

func (b *bboltStorage) Delete(id string) error {
	return nil
}

func (b *bboltStorage) Update(job models.Job) error {
	return nil
}*/
