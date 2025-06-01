package repo

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// SQLiteDB implements the DB interface for SQLite
type SQLiteDB struct {
	db     *sql.DB
	closed bool
}

// NewSQLiteDB creates a new SQLite database connection
func NewSQLiteDB(dbPath string) (*SQLiteDB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Initialize the database schema
	if err := initSchema(db); err != nil {
		return nil, err
	}

	return &SQLiteDB{
		db: db,
	}, nil
}

// Exec executes a query without returning any rows
func (s *SQLiteDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	if s.closed {
		return nil, sql.ErrConnDone
	}

	return s.db.Exec(query, args...)
}

// Query executes a query that returns rows
func (s *SQLiteDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if s.closed {
		return nil, sql.ErrConnDone
	}

	return s.db.Query(query, args...)
}

// QueryRow executes a query that returns a single row
func (s *SQLiteDB) QueryRow(query string, args ...interface{}) *sql.Row {
	if s.closed {
		return nil
	}

	return s.db.QueryRow(query, args...)
}

// Ping checks if the database connection is alive
func (s *SQLiteDB) Ping() error {
	if s.closed {
		return sql.ErrConnDone
	}

	return s.db.Ping()
}

// Close closes the database connection
func (s *SQLiteDB) Close() error {
	if s.closed {
		return nil
	}

	s.closed = true
	return s.db.Close()
}

// initSchema initializes the database schema if it doesn't exist
func initSchema(db *sql.DB) error {
	// Create novels table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS novels (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			original_title TEXT,
			cover TEXT,
			source TEXT NOT NULL,
			url TEXT NOT NULL,
			summary TEXT NOT NULL,
			author TEXT,
			status TEXT,
			genres TEXT,  -- Stored as JSON string
			chapters_count INTEGER NOT NULL DEFAULT 0,
			last_read_chapter_number INTEGER,
			last_read_timestamp INTEGER,
			last_updated INTEGER NOT NULL,
			date_added INTEGER NOT NULL
		);
	`)
	if err != nil {
		return err
	}

	// Create chapters table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS chapters (
			id TEXT PRIMARY KEY,
			novel_id TEXT NOT NULL,
			number INTEGER NOT NULL,
			title TEXT NOT NULL,
			original_title TEXT,
			content TEXT NOT NULL,
			date_translated INTEGER NOT NULL,
			word_count INTEGER,
			url TEXT,               
			FOREIGN KEY (novel_id) REFERENCES novels(id)
		);
	`)
	if err != nil {
		return err
	}

	// Create indexes for better performance
	_, err = db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_chapters_novel_id ON chapters(novel_id);
		CREATE INDEX IF NOT EXISTS idx_chapters_number ON chapters(novel_id, number);
	`)
	if err != nil {
		return err
	}

	return nil
}
