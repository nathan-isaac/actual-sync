package core

type Token = string

type TokenStore interface {
	FirstToken() (Token, error)
	AddToken(token Token) error
	HasToken(token Token) (bool, error)
}
