package repo

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// SQLiteDB implements the DB interface for SQLite
type SQLiteDB struct {
	db *sql.DB
}

// pragmas applied to every connection in the pool.
//
//	journal_mode=WAL   readers do not block the writer, which matters because a
//	                   chapter translation holds a write open for a long time
//	busy_timeout       wait rather than failing instantly on a locked database
//	foreign_keys       the chapters table declares one; SQLite ignores it unless
//	                   this is on, so deleting a novel could orphan chapters
//	synchronous=NORMAL safe with WAL and much faster than FULL
const dsnParams = "?_journal_mode=WAL&_busy_timeout=5000&_foreign_keys=on&_synchronous=NORMAL"

// NewSQLiteDB creates a new SQLite database connection
func NewSQLiteDB(dbPath string) (*SQLiteDB, error) {
	db, err := sql.Open("sqlite3", dbPath+dsnParams)
	if err != nil {
		return nil, err
	}

	// SQLite takes a single writer. Allowing the pool to open many connections
	// turns concurrent writes into "database is locked" errors instead of
	// queueing them.
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("open database %s: %w", dbPath, err)
	}

	if err = initSchema(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("initialise schema: %w", err)
	}

	return &SQLiteDB{db: db}, nil
}

// Exec executes a query without returning any rows
func (s *SQLiteDB) Exec(query string, args ...any) (sql.Result, error) {
	return s.db.Exec(query, args...)
}

// Query executes a query that returns rows
func (s *SQLiteDB) Query(query string, args ...any) (*sql.Rows, error) {
	return s.db.Query(query, args...)
}

// QueryRow executes a query that returns a single row.
//
// database/sql already reports use-after-close through the returned Row, so
// there is no guard here: the previous version returned a nil *sql.Row after
// Close, which panicked in every caller the moment they called Scan.
func (s *SQLiteDB) QueryRow(query string, args ...any) *sql.Row {
	return s.db.QueryRow(query, args...)
}

// Ping checks if the database connection is alive
func (s *SQLiteDB) Ping() error {
	return s.db.Ping()
}

// Close closes the database connection
func (s *SQLiteDB) Close() error {
	return s.db.Close()
}

// initSchema initializes the database schema if it doesn't exist
func initSchema(db *sql.DB) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS novels (
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
		)`,

		`CREATE TABLE IF NOT EXISTS chapters (
			id TEXT PRIMARY KEY,
			novel_id TEXT NOT NULL,
			number INTEGER NOT NULL,
			title TEXT NOT NULL,
			original_title TEXT,
			content TEXT NOT NULL,
			date_translated INTEGER NOT NULL,
			word_count INTEGER,
			url TEXT,
			next_chapter_url TEXT,
			FOREIGN KEY (novel_id) REFERENCES novels(id)
		)`,

		// Chapter lookups are always scoped to a novel.
		`CREATE INDEX IF NOT EXISTS idx_chapters_novel_id ON chapters(novel_id)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_chapters_novel_number ON chapters(novel_id, number)`,
		// Partial index: many chapters legitimately have no url, and a plain
		// UNIQUE index would reject the second empty string.
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_chapters_url
			ON chapters(url) WHERE url IS NOT NULL AND url <> ''`,

		// Every ORDER BY the novel list offers, plus the columns it filters on.
		`CREATE INDEX IF NOT EXISTS idx_novels_source ON novels(source)`,
		`CREATE INDEX IF NOT EXISTS idx_novels_status ON novels(status COLLATE NOCASE)`,
		`CREATE INDEX IF NOT EXISTS idx_novels_last_read_timestamp ON novels(last_read_timestamp DESC)`,
		`CREATE INDEX IF NOT EXISTS idx_novels_last_updated ON novels(last_updated DESC)`,
		`CREATE INDEX IF NOT EXISTS idx_novels_date_added ON novels(date_added DESC)`,
		`CREATE INDEX IF NOT EXISTS idx_novels_chapters_count ON novels(chapters_count DESC)`,
		`CREATE INDEX IF NOT EXISTS idx_novels_title ON novels(title COLLATE NOCASE)`,
		// Language shelves filter by source then order by recency.
		`CREATE INDEX IF NOT EXISTS idx_novels_source_last_read
			ON novels(source, last_read_timestamp DESC)`,
	}

	for _, statement := range statements {
		if _, err := db.Exec(statement); err != nil {
			return fmt.Errorf("%s: %w", statement, err)
		}
	}

	// Keeps the query planner's statistics current for the indexes above.
	if _, err := db.Exec(`ANALYZE`); err != nil {
		return fmt.Errorf("analyze: %w", err)
	}

	return nil
}
