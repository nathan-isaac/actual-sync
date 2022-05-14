package memory

import "github.com/nathanjisaac/actual-server-go/internal/core"

type PasswordStore struct {
	Passwords []core.Password
}

func New() core.PasswordStore {
	return &PasswordStore{Passwords: []core.Password{}}
}

func (it *PasswordStore) Count() (int, error) {
	return len(it.Passwords), nil
}

func (it *PasswordStore) First() (core.Password, error) {
	//TODO implement me
	panic("implement me")
}

func (it *PasswordStore) Add(password core.Password) error {
	it.Passwords = append(it.Passwords, password)

	return nil
}

func (it *PasswordStore) Set(password core.Password) error {
	//TODO implement me
	panic("implement me")
}
