package sqlite

import (
	"github.com/nathanjisaac/actual-server-go/internal/core"
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
	row, err := a.connection.First("SELECT * FROM sessions")

	if err != nil {
		return "", err
	}

	var password core.Token

	if err = row.Scan(&password); err != nil {
		return "", err
	}

	return password, nil
}

func (a *TokenStore) Add(token core.Token) error {
	_, _, err := a.connection.Mutate("INSERT INTO sessions (token) VALUES (?)", token)
	if err != nil {
		return err
	}

	return nil
}

func (a *TokenStore) Has(token core.Token) (bool, error) {
	rows, err := a.connection.All("SELECT * FROM sessions WHERE token = ?", token)
	if err != nil {
		return false, err
	}
	if rows.Next() {
		return true, nil
	} else {
		err = rows.Err()
		if err != nil {
			return false, err
		}
		return false, nil
	}
}
