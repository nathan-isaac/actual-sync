package memory

import (
	"github.com/nathanjisaac/actual-server-go/internal/core"
)

type TokenStore struct {
	Tokens []core.Token
}

func NewTokenStore() *TokenStore {
	return &TokenStore{
		Tokens: []core.Token{},
	}
}

func (a *TokenStore) First() (core.Token, error) {
	return a.Tokens[0], nil
}

func (a *TokenStore) Add(token core.Token) error {
	a.Tokens = append(a.Tokens, token)
	return nil
}

func (a *TokenStore) Has(token core.Token) (bool, error) {
	for _, v := range a.Tokens {
		if v == token {
			return true, nil
		}
	}
	return false, nil
}
