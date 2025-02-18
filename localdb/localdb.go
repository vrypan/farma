package localdb

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/vrypan/farma/config"
)

var db *badger.DB
var db_path = ""
var (
	ERR_NOT_FOUND  = errors.New("Not Found")
	ERR_NOT_STORED = errors.New("Not Stored")
)

//const dot_dir = ".fargo"

var FrameIdSequence *Sequence

func init() {
	if db_path == "" {
		configDir, err := config.ConfigDir()
		if err != nil {
			panic(err)
		}
		db_path = filepath.Join(configDir, "badger.db")
	}
}
func IsOpen() bool {
	return db != nil
}

func AssertOpen() {
	if db == nil {
		panic("DB not open")
	}
}

func Set(k []byte, v []byte) error {
	err := db.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry(k, v)
		return txn.SetEntry(e)
	})
	return err
}

// Write a key that will auto-expire after ttl*hours
func SetWithTtl(k string, v []byte, ttl int) error {
	err := db.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry([]byte(k), v).WithTTL(time.Duration(ttl) * time.Hour)
		return txn.SetEntry(e)
	})
	return err
}

func Get(k []byte) ([]byte, error) {
	var val []byte
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(k)
		if err != nil {
			return ERR_NOT_FOUND
		}
		err = item.Value(func(v []byte) error {
			val = v
			return nil
		})
		return err
	})
	if err != nil {
		return nil, err
	}
	return val, nil
}

func Delete(key []byte) error {
	return db.Update(func(txn *badger.Txn) error {
		err := txn.Delete(key)
		if err == badger.ErrKeyNotFound {
			return nil // No error if key doesn't exist
		}
		return err
	})
}

func Open() error {
	config.Load()
	var err error
	db, err = badger.Open(badger.DefaultOptions(db_path).WithLoggingLevel(badger.ERROR))
	if err != nil {
		return err
	}
	FrameIdSequence = NewSequence("FrameId", 5)
	if FrameIdSequence == nil {
		return fmt.Errorf("Unable to initialize FrameIdSequence")
	}
	return nil
}

func Close() error {
	FrameIdSequence.seq.Release()
	return db.Close()
}

func GetSize() (int64, error) {
	var size int64
	err := filepath.Walk(db_path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size, err
}

func CountEntries() (int, error) {
	AssertOpen()
	count := 0
	err := db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			count++
		}
		return nil
	})
	return count, err
}

func Path() string {
	return db_path
}

var Update = db.Update
var View = db.View

func GetKeys(prefix []byte, limit int) ([][]byte, error) {
	var keys [][]byte

	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()
		count := 0

		for it.Seek(prefix); it.ValidForPrefix(prefix) && count < limit; it.Next() {
			item := it.Item()
			keys = append(keys, item.Key())
			count++
		}
		return nil
	})

	return keys, err
}

// Each item in items[] is the value of the keys that match the prefix.
// lastKey is the last key that was returned, which can be used as the next cursor.
func GetPrefixP(prefix []byte, startKey []byte, limit int) (items [][]byte, lastKey []byte, err error) {
	err = db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		count := 0
		var v []byte
		for it.Seek(startKey); it.ValidForPrefix(prefix) && count < limit; it.Next() {
			item := it.Item()
			k := item.Key()
			v, err = item.ValueCopy(nil)
			if err != nil {
				return err
			}

			items = append(items, v)
			lastKey = k
			count++
		}
		return nil
	})

	// Return keys and the last key as the next cursor
	return items, lastKey, err
}

// Return the keys that match the prefix.
// lastKey is the last key that was returned, which can be used as the next cursor.
func GetKeysWithPrefix(prefix []byte, startKey []byte, limit int) (items [][]byte, lastKey []byte, err error) {
	err = db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		count := 0
		for it.Seek(startKey); it.ValidForPrefix(prefix) && count < limit; it.Next() {
			item := it.Item()
			k := item.Key()
			items = append(items, k)
			lastKey = k
			count++
		}
		return nil
	})

	// Return keys and the last key as the next cursor
	return items, lastKey, err
}
