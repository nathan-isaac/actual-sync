package sqlite

import (
	"database/sql"
)

type Connection struct {
	db *sql.DB
}

func NewConnection(filePath string) (*Connection, error) {
	// Create sql.DB and pass to connection
	// Run migrations if needed.
	return &Connection{}, nil
}

func (it *Connection) All(sql string, params ...any) (*sql.Rows, error) {
	//TODO implement me
	panic("implement me")
}

func (it *Connection) First(sql string, params ...any) (*sql.Row, error) {
	//TODO implement me
	panic("implement me")
}

// ...OTHER FUNCTIONS
