package storage

import (
	"encoding/json"
	"fmt"
	"time"

	"go.etcd.io/bbolt"
)

// BoltDB storage implementasyonu
type BoltStorage struct {
	db *bbolt.DB
}

// Yeni BoltDB storage oluşturur
func NewBoltStorage(dbPath string) (*BoltStorage, error) {
	db, err := bbolt.Open(dbPath, 0600, &bbolt.Options{
		Timeout: 1 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("BoltDB açılamadı: %v", err)
	}

	storage := &BoltStorage{db: db}
	if err := storage.createBuckets(); err != nil {
		return nil, fmt.Errorf("bucket'lar oluşturulamadı: %v", err)
	}

	return storage, nil
}

func (b *BoltStorage) createBuckets() error {
	return b.db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("urls"))
		return err
	})
}

func (b *BoltStorage) Save(url *URL) error {
	if url.CreatedAt.IsZero() {
		url.CreatedAt = time.Now()
	}

	return b.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("urls"))

		// ID üret (basit counter)
		if url.ID == 0 {
			id, _ := bucket.NextSequence()
			url.ID = int(id)
		}

		// JSON'a serialize et
		data, err := json.Marshal(url)
		if err != nil {
			return fmt.Errorf("serialization hatası: %v", err)
		}

		// Kaydet
		return bucket.Put([]byte(url.Short), data)
	})
}

func (b *BoltStorage) Get(shortCode string) (*URL, error) {
	var url URL

	err := b.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("urls"))
		data := bucket.Get([]byte(shortCode))

		if data == nil {
			return fmt.Errorf("URL bulunamadı")
		}

		return json.Unmarshal(data, &url)
	})

	if err != nil {
		return nil, err
	}
	return &url, nil
}

func (b *BoltStorage) IncrementClicks(shortCode string) error {
	return b.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("urls"))
		data := bucket.Get([]byte(shortCode))

		if data == nil {
			return fmt.Errorf("URL bulunamadı")
		}

		var url URL
		if err := json.Unmarshal(data, &url); err != nil {
			return err
		}

		url.Clicks++

		updatedData, err := json.Marshal(url)
		if err != nil {
			return err
		}

		return bucket.Put([]byte(shortCode), updatedData)
	})
}

func (b *BoltStorage) GetAll() ([]*URL, error) {
	var urls []*URL

	err := b.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("urls"))

		return bucket.ForEach(func(k, v []byte) error {
			var url URL
			if err := json.Unmarshal(v, &url); err != nil {
				return err
			}
			urls = append(urls, &url)
			return nil
		})
	})

	return urls, err
}

func (b *BoltStorage) Close() error {
	return b.db.Close()
}
