package sqlite

import (
	"github.com/nathanjisaac/actual-server-go/internal/storage"
	"github.com/stretchr/testify/assert"
	"testing"
)

func newStore(t *testing.T) *TokenStore {
	conn, err := NewConnection(":memory:")
	assert.NoError(t, err)

	return NewTokenStore(conn)
}

func TestTokenStore_First(t *testing.T) {
	t.Run("given no rows", func(t *testing.T) {
		store := newStore(t)

		_, err := store.First()

		assert.ErrorIs(t, err, storage.RecordNotFound)
	})

	t.Run("given one row", func(t *testing.T) {
		store := newStore(t)

		err := store.Add("token")
		assert.NoError(t, err)

		token, err := store.First()

		assert.NoError(t, err)
		assert.Equal(t, "token", token)
	})

	t.Run("given two rows then return first", func(t *testing.T) {
		store := newStore(t)

		err := store.Add("a")
		assert.NoError(t, err)
		err = store.Add("b")
		assert.NoError(t, err)

		token, err := store.First()

		assert.NoError(t, err)
		assert.Equal(t, "a", token)
	})
}

func TestTokenStore_Has(t *testing.T) {
	t.Run("given no rows", func(t *testing.T) {
		store := newStore(t)

		hasToken, err := store.Has("token")

		assert.NoError(t, err)
		assert.Equal(t, false, hasToken)
	})

	t.Run("given one row", func(t *testing.T) {
		store := newStore(t)

		err := store.Add("token")
		assert.NoError(t, err)

		hasToken, err := store.Has("token")

		assert.NoError(t, err)
		assert.Equal(t, true, hasToken)
	})

	t.Run("given one row with miss matching token", func(t *testing.T) {
		store := newStore(t)

		err := store.Add("token")
		assert.NoError(t, err)

		hasToken, err := store.Has("other")

		assert.NoError(t, err)
		assert.Equal(t, false, hasToken)
	})
}
