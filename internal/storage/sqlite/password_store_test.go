//nolint: dupl // Disabling dupl for tests. It detects similar testcases for different tests.
package sqlite_test

import (
	"testing"

	internal_errors "github.com/nathanjisaac/actual-server-go/internal/errors"
	"github.com/nathanjisaac/actual-server-go/internal/storage/sqlite"
	"github.com/stretchr/testify/assert"
)

func newTestPasswordStore(t *testing.T) (*sqlite.PasswordStore, *sqlite.Connection) {
	conn, err := sqlite.NewAccountConnection(":memory:")
	assert.NoError(t, err)

	return sqlite.NewPasswordStore(conn), conn
}

func TestPasswordStore_Count(t *testing.T) {
	t.Run("given no rows", func(t *testing.T) {
		store, conn := newTestPasswordStore(t)
		defer conn.Close()

		c, err := store.Count()

		assert.NoError(t, err)
		assert.Equal(t, 0, c)
	})

	t.Run("given two row", func(t *testing.T) {
		store, conn := newTestPasswordStore(t)
		defer conn.Close()

		err := store.Add("password0")
		assert.NoError(t, err)
		err = store.Add("password1")
		assert.NoError(t, err)

		c, err := store.Count()

		assert.NoError(t, err)
		assert.Equal(t, 2, c)
	})
}

func TestPasswordStore_First(t *testing.T) {
	t.Run("given no rows", func(t *testing.T) {
		store, conn := newTestPasswordStore(t)
		defer conn.Close()

		_, err := store.First()

		assert.ErrorIs(t, err, internal_errors.ErrStorageRecordNotFound)
	})

	t.Run("given one row", func(t *testing.T) {
		store, conn := newTestPasswordStore(t)
		defer conn.Close()

		err := store.Add("password")
		assert.NoError(t, err)

		password, err := store.First()

		assert.NoError(t, err)
		assert.Equal(t, "password", password)
	})

	t.Run("given two rows then return first", func(t *testing.T) {
		store, conn := newTestPasswordStore(t)
		defer conn.Close()

		err := store.Add("a")
		assert.NoError(t, err)
		err = store.Add("b")
		assert.NoError(t, err)

		password, err := store.First()

		assert.NoError(t, err)
		assert.Equal(t, "a", password)
	})
}

func TestPasswordStore_Add(t *testing.T) {
	t.Run("add new row", func(t *testing.T) {
		store, conn := newTestPasswordStore(t)
		defer conn.Close()

		err := store.Add("password")
		assert.NoError(t, err)

		p, err := store.First()

		assert.NoError(t, err)
		assert.Equal(t, "password", p)
	})
}

func TestPasswordStore_Set(t *testing.T) {
	t.Run("given no row", func(t *testing.T) {
		store, conn := newTestPasswordStore(t)
		defer conn.Close()

		err := store.Set("password")
		assert.ErrorIs(t, err, internal_errors.ErrStorageNoRecordUpdated)
	})

	t.Run("given one row", func(t *testing.T) {
		store, conn := newTestPasswordStore(t)
		defer conn.Close()

		err := store.Add("password")
		assert.NoError(t, err)

		p, err := store.First()

		assert.NoError(t, err)
		assert.Equal(t, "password", p)

		err = store.Set("newPassword")
		assert.NoError(t, err)

		p, err = store.First()

		assert.NoError(t, err)
		assert.Equal(t, "newPassword", p)
	})
}
