package sqlite

import (
	"database/sql"

	"github.com/nathanjisaac/actual-server-go/internal/core"
	"github.com/nathanjisaac/actual-server-go/internal/errors"
)

type TokenStore struct {
	connection *Connection
}

func NewTokenStore(connection *Connection) *TokenStore {
	return &TokenStore{
		connection: connection,
	}
}

func (a *TokenStore) First() (core.Token, error) {
	var token core.Token

	row, err := a.connection.First("SELECT * FROM sessions")

	if err != nil {
		return token, err
	}

	if err = row.Scan(&token); err != nil {
		if err == sql.ErrNoRows {
			return token, errors.StorageErrorRecordNotFound
		}
		return token, err
	}

	return token, nil
}

func (a *TokenStore) Add(token core.Token) error {
	_, _, err := a.connection.Mutate("INSERT INTO sessions (token) VALUES (?)", token)
	if err != nil {
		return err
	}

	return nil
}

func (a *TokenStore) Has(token core.Token) (bool, error) {
	var count int

	row, err := a.connection.First("SELECT count(*) FROM sessions WHERE token = ?", token)
	if err != nil {
		return false, err
	}

	if err = row.Scan(&count); err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}
