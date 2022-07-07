package sqlite

import (
	"database/sql"

	"github.com/nathanjisaac/actual-server-go/internal/core"
	"github.com/nathanjisaac/actual-server-go/internal/errors"
)

type PasswordStore struct {
	connection *Connection
}

func NewPasswordStore(connection *Connection) *PasswordStore {
	return &PasswordStore{
		connection: connection,
	}
}

func (a *PasswordStore) Count() (int, error) {
	row, err := a.connection.First("SELECT count(*) FROM auth")

	if err != nil {
		return 0, err
	}

	var count int

	if err = row.Scan(&count); err != nil {
		return 0, errors.StorageErrorRecordNotFound
	}

	return count, nil
}

func (a *PasswordStore) First() (core.Password, error) {
	var password core.Password

	row, err := a.connection.First("SELECT * FROM auth")

	if err != nil {
		return password, err
	}

	if err = row.Scan(&password); err != nil {
		if err == sql.ErrNoRows {
			return password, errors.StorageErrorRecordNotFound
		}
		return password, err
	}

	return password, nil
}

func (a *PasswordStore) Add(password core.Password) error {
	_, _, err := a.connection.Mutate("INSERT INTO auth (password) VALUES (?)", password)
	if err != nil {
		return err
	}

	return nil
}

func (a *PasswordStore) Set(password core.Password) error {
	rows, _, err := a.connection.Mutate("UPDATE auth SET password = ?", password)
	if err != nil {
		return err
	} else if rows == 0 {
		return errors.StorageErrorNoRecordUpdated
	}

	return nil
}
