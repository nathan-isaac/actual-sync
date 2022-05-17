package sqlite_test

import (
	"testing"

	"github.com/nathanjisaac/actual-server-go/internal/core"
	"github.com/nathanjisaac/actual-server-go/internal/core/crdt/timestamp"
	"github.com/nathanjisaac/actual-server-go/internal/storage/sqlite"
	"github.com/stretchr/testify/assert"
)

func newTestMessageStore(t *testing.T) (*sqlite.MessageStore, *sqlite.Connection) {
	conn, err := sqlite.NewMessageConnection(":memory:")
	assert.NoError(t, err)

	return sqlite.NewMessageStore(conn), conn
}

func TestMessageStore_Add(t *testing.T) {
	t.Run("add one row", func(t *testing.T) {
		store, conn := newTestMessageStore(t)
		defer conn.Close()

		ts := timestamp.NewTimestamp(1000000000000, 5678, "ABCDEFGH12345678")
		msg := core.BinaryMessage{Timestamp: ts.ToString(), IsEncrypted: true, Content: []byte{11, 12, 13}}
		flag, err := store.Add(msg)

		assert.NoError(t, err)
		assert.Equal(t, true, flag)
	})

	t.Run("add second row with same primarykey", func(t *testing.T) {
		store, conn := newTestMessageStore(t)
		defer conn.Close()

		ts := timestamp.NewTimestamp(1000000000000, 5678, "ABCDEFGH12345678")
		msg := core.BinaryMessage{Timestamp: ts.ToString(), IsEncrypted: true, Content: []byte{11, 12, 13}}
		flag, err := store.Add(msg)

		assert.NoError(t, err)
		assert.Equal(t, true, flag)

		ts = timestamp.NewTimestamp(1000000000000, 5678, "ABCDEFGH12345678")
		msg = core.BinaryMessage{Timestamp: ts.ToString(), IsEncrypted: true, Content: []byte{19, 1, 23}}
		flag, err = store.Add(msg)

		assert.NoError(t, err)
		assert.Equal(t, false, flag)
	})
}

func TestMessageStore_GetSince(t *testing.T) {
	t.Run("given two rows and return two rows", func(t *testing.T) {
		store, conn := newTestMessageStore(t)
		defer conn.Close()

		ts1 := timestamp.NewTimestamp(1000000000000, 5678, "ABCDEFGH12345678")
		msg1 := core.BinaryMessage{Timestamp: ts1.ToString(), IsEncrypted: true, Content: []byte{11, 12, 13}}
		flag, err := store.Add(msg1)

		assert.NoError(t, err)
		assert.Equal(t, true, flag)

		ts2 := timestamp.NewTimestamp(1000, 1234, "12345678ABCDEFGH")
		msg2 := core.BinaryMessage{Timestamp: ts2.ToString(), IsEncrypted: true, Content: []byte{1, 25, 17}}
		flag, err = store.Add(msg2)

		assert.NoError(t, err)
		assert.Equal(t, true, flag)

		since := timestamp.NewTimestamp(0, 5678, "ABCDEFGH12345678")
		messages, err := store.GetSince(since.ToString())

		assert.NoError(t, err)
		assert.Equal(t, 2, len(messages))
		assert.Equal(t, &core.BinaryMessage{Timestamp: ts1.ToString(), IsEncrypted: true, Content: []byte{11, 12, 13}}, messages[1])
		assert.Equal(t, &core.BinaryMessage{Timestamp: ts2.ToString(), IsEncrypted: true, Content: []byte{1, 25, 17}}, messages[0])
	})

	t.Run("given 2 rows and get rows since inbetween", func(t *testing.T) {
		store, conn := newTestMessageStore(t)
		defer conn.Close()

		ts1 := timestamp.NewTimestamp(1000000000000, 5678, "ABCDEFGH12345678")
		msg1 := core.BinaryMessage{Timestamp: ts1.ToString(), IsEncrypted: true, Content: []byte{11, 12, 13}}
		flag, err := store.Add(msg1)

		assert.NoError(t, err)
		assert.Equal(t, true, flag)

		ts2 := timestamp.NewTimestamp(0, 1234, "12345678ABCDEFGH")
		msg2 := core.BinaryMessage{Timestamp: ts2.ToString(), IsEncrypted: true, Content: []byte{1, 25, 17}}
		flag, err = store.Add(msg2)

		assert.NoError(t, err)
		assert.Equal(t, true, flag)

		since := timestamp.NewTimestamp(100000000000, 5678, "ABCDEFGH12345678")
		messages, err := store.GetSince(since.ToString())

		assert.NoError(t, err)
		assert.Equal(t, 1, len(messages))
		assert.Equal(t, &core.BinaryMessage{Timestamp: ts1.ToString(), IsEncrypted: true, Content: []byte{11, 12, 13}}, messages[0])
	})

	t.Run("given 2 rows and returns no rows", func(t *testing.T) {
		store, conn := newTestMessageStore(t)
		defer conn.Close()

		ts := timestamp.NewTimestamp(1000000000000, 5678, "ABCDEFGH12345678")
		msg := core.BinaryMessage{Timestamp: ts.ToString(), IsEncrypted: true, Content: []byte{11, 12, 13}}
		flag, err := store.Add(msg)

		assert.NoError(t, err)
		assert.Equal(t, true, flag)

		ts = timestamp.NewTimestamp(0, 1234, "12345678ABCDEFGH")
		msg = core.BinaryMessage{Timestamp: ts.ToString(), IsEncrypted: true, Content: []byte{1, 25, 17}}
		flag, err = store.Add(msg)

		assert.NoError(t, err)
		assert.Equal(t, true, flag)

		since := timestamp.NewTimestamp(10000000000000, 5678, "ABCDEFGH12345678")
		messages, err := store.GetSince(since.ToString())

		assert.NoError(t, err)
		assert.Equal(t, 0, len(messages))
	})
}
