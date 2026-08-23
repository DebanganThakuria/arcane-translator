package repo

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// DB is an interface for database operations
type DB interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	Ping() error
	Close() error
}
