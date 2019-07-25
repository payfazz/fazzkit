package message

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"

	bolt "go.etcd.io/bbolt"
)

const bucketName = "event_message"

//BoltRepo implements repo interface
type BoltRepo struct {
	db       *bolt.DB
	name     string
	isClosed bool
}

var mutex sync.Mutex
var repos = map[string]*BoltRepo{}

//NewBoltRepo initialize bolt repository, open bolt connection to file
func NewBoltRepo(dbName string) *BoltRepo {
	mutex.Lock()
	defer mutex.Unlock()

	folderPath := os.Getenv("TMP_PATH")
	dbName = dbName + ".boltdb"

	if r, ok := repos[dbName]; ok {
		if !r.isClosed {
			return r
		}

		r.db, _ = bolt.Open(filepath.Join(folderPath, dbName), 0666, &bolt.Options{Timeout: 5 * time.Second})
		r.isClosed = false
		return r
	}

	db, _ := bolt.Open(filepath.Join(folderPath, dbName), 0666, &bolt.Options{Timeout: 5 * time.Second})
	db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte(bucketName))
		return nil
	})

	repos[dbName] = &BoltRepo{db, dbName, false}
	return repos[dbName]
}

//TotalBoltRepo get current total bolt repository in flyweight
func TotalBoltRepo() int {
	return len(repos)
}

//Save message to bolt db
func (repo *BoltRepo) Save(message Message) error {
	if repo.isClosed {
		return errors.New("db closed")
	}

	repo.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		err := b.Put(message.ID, message.Data)
		return err
	})
	return nil
}

//Load message from bolt db
func (repo *BoltRepo) Load(ID []byte) (*Message, error) {
	if repo.isClosed {
		return nil, errors.New("db closed")
	}

	message := &Message{}

	repo.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		data := b.Get(ID)
		if data == nil {
			message = nil
			return nil
		}

		message = &Message{
			ID:   ID,
			Data: data,
		}

		return nil
	})

	return message, nil
}

//Close close db connection
func (repo *BoltRepo) Close() {
	if repo.db != nil {
		repo.db.Close()
		repo.isClosed = true
	}
}
