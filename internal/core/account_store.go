package core

type Password = string
type Token = string

type AccountStore interface {
	PasswordsCount() (int, error)
	FirstPassword() (Password, error)
	FirstToken() (Token, error)
	AddPassword(password Password) error
	SetPassword(password Password) error
	AddToken(token Token) error
	HasToken(token Token) (bool, error)
}
