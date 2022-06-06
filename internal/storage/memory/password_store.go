package memory

import (
	"github.com/nathanjisaac/actual-server-go/internal/core"
	"github.com/nathanjisaac/actual-server-go/internal/storage"
)

type PasswordStore struct {
	Passwords []core.Password
}

func NewPasswordStore() *PasswordStore {
	return &PasswordStore{
		Passwords: []core.Password{},
	}
}

func (it *PasswordStore) Count() (int, error) {
	return len(it.Passwords), nil
}

func (it *PasswordStore) First() (core.Password, error) {
	if len(it.Passwords) == 0 {
		return "", storage.ErrorRecordNotFound
	}
	return it.Passwords[len(it.Passwords)-1], nil
}

func (it *PasswordStore) Add(password core.Password) error {
	it.Passwords = append(it.Passwords, password)
	return nil
}

func (it *PasswordStore) Set(password core.Password) error {
	it.Passwords[len(it.Passwords)-1] = password
	return nil
}
