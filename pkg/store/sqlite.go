package store

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type SQLiteStore struct {
	db         *sql.DB
	maxRecords int
}

func NewSQLiteStore(path string, maxRecords int) (*SQLiteStore, error) {
	if dir := filepath.Dir(path); dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, err
		}
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS tokens (
		id    INTEGER PRIMARY KEY AUTOINCREMENT,
		token TEXT    UNIQUE NOT NULL
	)`)
	if err != nil {
		db.Close()
		return nil, err
	}

	return &SQLiteStore{db: db, maxRecords: maxRecords}, nil
}

func (s *SQLiteStore) Exists(token string) (bool, error) {
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM tokens WHERE token = ?)", token).Scan(&exists)
	return exists, err
}

func (s *SQLiteStore) Add(token string) error {
	_, err := s.db.Exec("INSERT OR IGNORE INTO tokens (token) VALUES (?)", token)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(`DELETE FROM tokens WHERE id NOT IN (
		SELECT id FROM tokens ORDER BY id DESC LIMIT ?
	)`, s.maxRecords)
	return err
}

func (s *SQLiteStore) Close() error { return s.db.Close() }
