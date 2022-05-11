package sqlite

import (
	"github.com/nathanjisaac/actual-server-go/internal/core"
)

type AccountStore struct {
	connection *Connection
}

func NewAccountStore(connection *Connection) *AccountStore {
	return &AccountStore{
		connection: connection,
	}
}

func (a *AccountStore) PasswordsCount() (int, error) {
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

func (a *AccountStore) FirstPassword() (core.Password, error) {
	//TODO implement me
	panic("implement me")
}

func (a *AccountStore) FirstToken() (core.Token, error) {
	//TODO implement me
	panic("implement me")
}

func (a *AccountStore) AddPassword(password core.Password) error {
	//TODO implement me
	panic("implement me")
}

func (a *AccountStore) SetPassword(password core.Password) error {
	//TODO implement me
	panic("implement me")
}

func (a *AccountStore) AddToken(token core.Token) error {
	//TODO implement me
	panic("implement me")
}

func (a *AccountStore) HasToken(token core.Token) (bool, error) {
	//TODO implement me
	panic("implement me")
}
