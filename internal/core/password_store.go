package core

type Password = string

type PasswordStore interface {
	PasswordsCount() (int, error)
	FirstPassword() (Password, error)
	AddPassword(password Password) error
	SetPassword(password Password) error
}
