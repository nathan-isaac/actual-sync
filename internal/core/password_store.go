package core

type Password = string

type PasswordStore interface {
	Count() (int, error)
	First() (Password, error)
	Add(password Password) error
	Set(password Password) error
}
