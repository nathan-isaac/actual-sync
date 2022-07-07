package core

import "database/sql"

type StorageType string

const (
	Sqlite StorageType = "sqlite"
)

type StorageConfig interface {
}

type Connection interface {
	All(string, ...any) (*sql.Rows, error)
	First(string, ...any) (*sql.Row, error)
	Exec(string) (sql.Result, error)
	Mutate(string, ...any) (int64, int64, error)
	Transaction(func(*sql.Tx) error) error
	Close() error
}
