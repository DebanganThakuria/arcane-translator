package repo

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// DB is an interface for database operations
type DB interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Ping() error
	Close() error
}
