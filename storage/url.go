package storage

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// URL yapısı - kısaltılmış URL'leri saklamak için
type URL struct {
	ID        int       `json:"id"`         // Kısa URL'nin ID'si
	Original  string    `json:"original"`   // Orijinal URL
	Short     string    `json:"short"`      // Kısaltılmış URL
	CreatedAt time.Time `json:"created_at"` // Oluşturulma zamanı
	Clicks    int       `json:"clicks"`     // Tıklanma sayısı
}

// Storage interface - farklı storage tiplerini desteklemek için
type Storage interface {
	Save(url *URL) error
	Get(shortCode string) (*URL, error)
	IncrementClicks(shortCode string) error
	GetAll() ([]*URL, error)
	Close() error
}

// In-memory storage implementasyonu
type MemoryStorage struct {
	urls    map[string]*URL
	counter int
	mutex   sync.RWMutex
}

// SQLite storage implementasyonu
type SQLiteStorage struct {
	db *sql.DB
}

// Memory storage oluşturur
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		urls:    make(map[string]*URL),
		counter: 1,
	}
}

// SQLite storage oluşturur
func NewSQLiteStorage(dbPath string) (*SQLiteStorage, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("SQLite veritabanı açılamadı: %v", err)
	}

	storage := &SQLiteStorage{db: db}
	if err := storage.createTables(); err != nil {
		return nil, fmt.Errorf("tablolar oluşturulamadı: %v", err)
	}

	return storage, nil
}

// MEMORY STORAGE METHODS

func (m *MemoryStorage) Save(url *URL) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if url.ID == 0 {
		url.ID = m.counter
		m.counter++
	}

	if url.CreatedAt.IsZero() {
		url.CreatedAt = time.Now()
	}

	m.urls[url.Short] = url
	return nil
}

func (m *MemoryStorage) Get(shortCode string) (*URL, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	url, exists := m.urls[shortCode]
	if !exists {
		return nil, fmt.Errorf("URL bulunamadı")
	}
	return url, nil
}

func (m *MemoryStorage) IncrementClicks(shortCode string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	url, exists := m.urls[shortCode]
	if !exists {
		return fmt.Errorf("URL bulunamadı")
	}

	url.Clicks++
	return nil
}

func (m *MemoryStorage) GetAll() ([]*URL, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	urls := make([]*URL, 0, len(m.urls))
	for _, url := range m.urls {
		urls = append(urls, url)
	}
	return urls, nil
}

func (m *MemoryStorage) Close() error {
	// Memory storage için kapatma işlemi yok
	return nil
}

// SQLITE STORAGE METHODS

func (s *SQLiteStorage) createTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS urls (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		original TEXT NOT NULL,
		short TEXT UNIQUE NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		clicks INTEGER DEFAULT 0
	);
	
	CREATE INDEX IF NOT EXISTS idx_short ON urls(short);
	`

	_, err := s.db.Exec(query)
	return err
}

func (s *SQLiteStorage) Save(url *URL) error {
	if url.CreatedAt.IsZero() {
		url.CreatedAt = time.Now()
	}

	query := `
	INSERT INTO urls (original, short, created_at, clicks) 
	VALUES (?, ?, ?, ?)
	`

	result, err := s.db.Exec(query, url.Original, url.Short, url.CreatedAt, url.Clicks)
	if err != nil {
		return fmt.Errorf("URL kaydedilemedi: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("ID alınamadı: %v", err)
	}

	url.ID = int(id)
	return nil
}

func (s *SQLiteStorage) Get(shortCode string) (*URL, error) {
	query := `
	SELECT id, original, short, created_at, clicks 
	FROM urls WHERE short = ?
	`

	row := s.db.QueryRow(query, shortCode)

	var url URL
	var createdAtStr string

	err := row.Scan(&url.ID, &url.Original, &url.Short, &createdAtStr, &url.Clicks)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("URL bulunamadı")
		}
		return nil, fmt.Errorf("URL getirilemedi: %v", err)
	}

	// Parse created_at
	url.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
	if err != nil {
		url.CreatedAt = time.Now() // Fallback
	}

	return &url, nil
}

func (s *SQLiteStorage) IncrementClicks(shortCode string) error {
	query := `UPDATE urls SET clicks = clicks + 1 WHERE short = ?`

	result, err := s.db.Exec(query, shortCode)
	if err != nil {
		return fmt.Errorf("tıklama sayısı güncellenemedi: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("URL bulunamadı")
	}

	return nil
}

func (s *SQLiteStorage) GetAll() ([]*URL, error) {
	query := `
	SELECT id, original, short, created_at, clicks 
	FROM urls ORDER BY created_at DESC
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("URL'ler getirilemedi: %v", err)
	}
	defer rows.Close()

	var urls []*URL

	for rows.Next() {
		var url URL
		var createdAtStr string

		err := rows.Scan(&url.ID, &url.Original, &url.Short, &createdAtStr, &url.Clicks)
		if err != nil {
			continue // Hatalı satırı atla
		}

		// Parse created_at
		url.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
		if err != nil {
			url.CreatedAt = time.Now() // Fallback
		}

		urls = append(urls, &url)
	}

	return urls, nil
}

func (s *SQLiteStorage) Close() error {
	return s.db.Close()
}

// FACTORY FUNCTION

// Config'e göre storage oluşturur
func NewStorage(dbType string) (Storage, error) {
	switch dbType {
	case "memory":
		return NewMemoryStorage(), nil
	case "sqlite":
		return NewSQLiteStorage("urls.db")
	case "bolt":
		return NewBoltStorage("urls.bolt.db")
	default:
		return nil, fmt.Errorf("desteklenmeyen veritabanı tipi: %s (desteklenen: memory, sqlite, bbolt)", dbType)
	}
}
