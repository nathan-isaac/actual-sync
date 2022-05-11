package sqlite

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

type Connection struct {
	db *sql.DB
}

func NewConnection(filePath string) (*Connection, error) {
	db, err := sql.Open("sqlite", filePath)
	if err != nil {
		return nil, err
	}

	// TODO: Run migrations if needed.

	return &Connection{db: db}, nil
}

func (it *Connection) All(sqlString string, params ...any) (*sql.Rows, error) {
	stmt, err := it.db.Prepare(sqlString)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query(params...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return rows, nil
}

func (it *Connection) First(sqlString string, params ...any) (*sql.Row, error) {
	stmt, err := it.db.Prepare(sqlString)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	row := stmt.QueryRow(params...)
	err = row.Err()
	if err != nil {
		return nil, err
	}

	return row, nil
}

func (it *Connection) Exec(sqlString string) (sql.Result, error) {
	return it.db.Exec(sqlString)
}

func (it *Connection) Mutate(sqlString string, params ...any) (int64, int64, error) {
	stmt, err := it.db.Prepare(sqlString)
	if err != nil {
		return 0, 0, err
	}

	defer stmt.Close()

	result, err := stmt.Exec(params...)
	if err != nil {
		return 0, 0, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return 0, 0, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, 0, err
	}

	return rows, lastId, nil
}

func (it *Connection) Transaction(fn func(*sql.Tx) error) error {
	tx, err := it.db.Begin()
	if err != nil {
		return err
	}

	err = fn(tx)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Fatalf("DB transaction FAILURE: unable to rollback: %v", rollbackErr)
		}
		log.Fatal(err)
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (it *Connection) Close() error {
	return it.db.Close()
}
