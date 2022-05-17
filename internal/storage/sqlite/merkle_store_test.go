package sqlite_test

import (
	"testing"

	"github.com/nathanjisaac/actual-server-go/internal/core"
	"github.com/nathanjisaac/actual-server-go/internal/storage"
	"github.com/nathanjisaac/actual-server-go/internal/storage/sqlite"
	"github.com/stretchr/testify/assert"
)

func newTestMerkleStore(t *testing.T) (*sqlite.MerkleStore, *sqlite.Connection) {
	conn, err := sqlite.NewMessageConnection(":memory:")
	assert.NoError(t, err)

	return sqlite.NewMerkleStore(conn), conn
}

func TestMerkleStore_Add(t *testing.T) {
	t.Run("add one row", func(t *testing.T) {
		store, conn := newTestMerkleStore(t)
		defer conn.Close()

		msg := core.MerkleMessage{MerkleId: "1", Merkle: "stringifiedMerkle"}
		err := store.Add(msg)

		assert.NoError(t, err)
	})

	t.Run("add second row with same primarykey", func(t *testing.T) {
		store, conn := newTestMerkleStore(t)
		defer conn.Close()

		msg := core.MerkleMessage{MerkleId: "1", Merkle: "stringifiedMerkle1"}
		err := store.Add(msg)

		assert.NoError(t, err)

		msg = core.MerkleMessage{MerkleId: "1", Merkle: "stringifiedMerkle2"}
		err = store.Add(msg)

		assert.NoError(t, err)
	})
}

func TestMerkleStore_GetForGroup(t *testing.T) {
	t.Run("given no rows", func(t *testing.T) {
		store, conn := newTestMerkleStore(t)
		defer conn.Close()

		_, err := store.GetForGroup("1")

		assert.ErrorIs(t, err, storage.ErrorRecordNotFound)
	})

	t.Run("given two rows and return one", func(t *testing.T) {
		store, conn := newTestMerkleStore(t)
		defer conn.Close()

		msg := core.MerkleMessage{MerkleId: "1", Merkle: "stringifiedMerkle1"}
		err := store.Add(msg)

		assert.NoError(t, err)

		msg = core.MerkleMessage{MerkleId: "2", Merkle: "stringifiedMerkle2"}
		err = store.Add(msg)

		assert.NoError(t, err)

		message, err := store.GetForGroup("1")

		assert.NoError(t, err)
		assert.Equal(t, &core.MerkleMessage{MerkleId: "1", Merkle: "stringifiedMerkle1"}, message)
	})
}
