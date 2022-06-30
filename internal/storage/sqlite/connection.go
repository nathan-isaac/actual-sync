package sqlite

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	migrateSqlite "github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "modernc.org/sqlite" // Using blank-import due to requirement by package
)

type Connection struct {
	db *sql.DB
}

//go:embed migrations/account/*.sql
var migrationsAccount embed.FS

//go:embed migrations/message/*.sql
var migrationsMessage embed.FS

func NewAccountConnection(dataSource string) (*Connection, error) {
	db, err := sql.Open("sqlite", dataSource)
	if err != nil {
		return nil, err
	}

	sourceDriver, err := iofs.New(migrationsAccount, "migrations/account")
	if err != nil {
		return nil, err
	}
	defer sourceDriver.Close()
	dbDriver, err := migrateSqlite.WithInstance(db, &migrateSqlite.Config{})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithInstance("iofs", sourceDriver, "sqlite", dbDriver)
	if err != nil {
		return nil, err
	}

	// Migrate to latest schema
	err = m.Up()
	if err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return nil, err
		}
	}

	conn := &Connection{db: db}

	return conn, nil
}

func NewMessageConnection(dataSource string) (*Connection, error) {
	db, err := sql.Open("sqlite", dataSource)
	if err != nil {
		return nil, err
	}

	sourceDriver, err := iofs.New(migrationsMessage, "migrations/message")
	if err != nil {
		return nil, err
	}
	defer sourceDriver.Close()
	dbDriver, err := migrateSqlite.WithInstance(db, &migrateSqlite.Config{})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithInstance("iofs", sourceDriver, "sqlite", dbDriver)
	if err != nil {
		return nil, err
	}

	// Migrate to latest schema
	err = m.Up()
	if err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return nil, err
		}
	}

	conn := &Connection{db: db}

	return conn, nil
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

	return rows, nil
}

func (it *Connection) First(sqlString string, params ...any) (*sql.Row, error) {
	stmt, err := it.db.Prepare(sqlString)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	row := stmt.QueryRow(params...)

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

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, 0, err
	}

	return rows, lastID, nil
}

func (it *Connection) Transaction(fn func(*sql.Tx) error) error {
	tx, err := it.db.Begin()
	if err != nil {
		return err
	}

	err = fn(tx)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("DB transaction FAILURE: unable to rollback: %w", rollbackErr)
		}
		return err
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
