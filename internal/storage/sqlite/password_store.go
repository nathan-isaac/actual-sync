package sqlite

import (
	"github.com/nathanjisaac/actual-server-go/internal/core"
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
		return 0, err
	}

	return count, nil
}

func (a *PasswordStore) First() (core.Password, error) {
	//TODO implement me
	panic("implement me")
}

func (a *PasswordStore) Add(password core.Password) error {
	//TODO implement me
	panic("implement me")
}

func (a *PasswordStore) Set(password core.Password) error {
	//TODO implement me
	panic("implement me")
}
