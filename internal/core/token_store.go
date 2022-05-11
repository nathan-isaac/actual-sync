package core

type Token = string

type TokenStore interface {
	First() (Token, error)
	Add(token Token) error
	Has(token Token) (bool, error)
}
