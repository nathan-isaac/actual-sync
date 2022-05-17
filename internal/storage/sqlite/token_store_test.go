package sqlite_test

import (
	"testing"

	"github.com/nathanjisaac/actual-server-go/internal/storage"
	"github.com/nathanjisaac/actual-server-go/internal/storage/sqlite"
	"github.com/stretchr/testify/assert"
)

func newTestTokenStore(t *testing.T) (*sqlite.TokenStore, *sqlite.Connection) {
	conn, err := sqlite.NewAccountConnection(":memory:")
	assert.NoError(t, err)

	return sqlite.NewTokenStore(conn), conn
}

func TestTokenStore_First(t *testing.T) {
	t.Run("given no rows", func(t *testing.T) {
		store, conn := newTestTokenStore(t)
		defer conn.Close()

		_, err := store.First()

		assert.ErrorIs(t, err, storage.ErrorRecordNotFound)
	})

	t.Run("given one row", func(t *testing.T) {
		store, conn := newTestTokenStore(t)
		defer conn.Close()

		err := store.Add("token")
		assert.NoError(t, err)

		token, err := store.First()

		assert.NoError(t, err)
		assert.Equal(t, "token", token)
	})

	t.Run("given two rows then return first", func(t *testing.T) {
		store, conn := newTestTokenStore(t)
		defer conn.Close()

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
		store, conn := newTestTokenStore(t)
		defer conn.Close()

		hasToken, err := store.Has("token")

		assert.NoError(t, err)
		assert.Equal(t, false, hasToken)
	})

	t.Run("given one row", func(t *testing.T) {
		store, conn := newTestTokenStore(t)
		defer conn.Close()

		err := store.Add("token")
		assert.NoError(t, err)

		hasToken, err := store.Has("token")

		assert.NoError(t, err)
		assert.Equal(t, true, hasToken)
	})

	t.Run("given one row with miss matching token", func(t *testing.T) {
		store, conn := newTestTokenStore(t)
		defer conn.Close()

		err := store.Add("token")
		assert.NoError(t, err)

		hasToken, err := store.Has("other")

		assert.NoError(t, err)
		assert.Equal(t, false, hasToken)
	})
}
