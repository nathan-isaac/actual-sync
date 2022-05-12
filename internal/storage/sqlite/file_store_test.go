package sqlite_test

import (
	"testing"

	"github.com/nathanjisaac/actual-server-go/internal/core"
	"github.com/nathanjisaac/actual-server-go/internal/storage"
	"github.com/nathanjisaac/actual-server-go/internal/storage/sqlite"
	"github.com/stretchr/testify/assert"
)

func newTestFileStore(t *testing.T) (*sqlite.FileStore, *sqlite.Connection) {
	conn, err := sqlite.NewConnection(":memory:")
	assert.NoError(t, err)

	return sqlite.NewFileStore(conn), conn
}

func TestFileStore_ForId(t *testing.T) {
	t.Run("given no rows", func(t *testing.T) {
		store, conn := newTestFileStore(t)
		defer conn.Close()

		_, err := store.ForId("1", false)

		assert.ErrorIs(t, err, storage.ErrorRecordNotFound)
	})

	t.Run("given three rows returns second", func(t *testing.T) {
		store, conn := newTestFileStore(t)
		defer conn.Close()

		err := store.Add(&core.NewFile{FileId: "1", GroupId: "g1", SyncVersion: 1, EncryptMeta: "A1B2C3", Name: "Budget"})
		assert.NoError(t, err)
		err = store.UpdateEncryption("1", "salt1", "keyid1", "test1")
		assert.NoError(t, err)

		err = store.Add(&core.NewFile{FileId: "2", GroupId: "g2", SyncVersion: 2, EncryptMeta: "B1F2G3", Name: "Budget2"})
		assert.NoError(t, err)
		err = store.UpdateEncryption("2", "salt2", "keyid2", "test2")
		assert.NoError(t, err)

		err = store.Add(&core.NewFile{FileId: "3", GroupId: "g3", SyncVersion: 3, EncryptMeta: "B4F7G9", Name: "Budget3"})
		assert.NoError(t, err)
		err = store.UpdateEncryption("3", "salt3", "keyid3", "test3")
		assert.NoError(t, err)

		f, err := store.ForId("2", false)

		assert.NoError(t, err)
		assert.Equal(t, &core.File{FileId: "2", GroupId: "g2", SyncVersion: 2, EncryptMeta: "B1F2G3", EncryptSalt: "salt2", EncryptKeyId: "keyid2", EncryptTest: "test2", Deleted: false, Name: "Budget2"}, f)
	})
}

func TestFileStore_All(t *testing.T) {
	t.Run("given no rows", func(t *testing.T) {
		store, conn := newTestFileStore(t)
		defer conn.Close()

		files, err := store.All()

		assert.NoError(t, err)
		assert.Equal(t, 0, len(files))
	})

	t.Run("given three rows returns all", func(t *testing.T) {
		store, conn := newTestFileStore(t)
		defer conn.Close()

		err := store.Add(&core.NewFile{FileId: "1", GroupId: "g1", SyncVersion: 1, EncryptMeta: "A1B2C3", Name: "Budget1"})
		assert.NoError(t, err)
		err = store.UpdateEncryption("1", "salt1", "keyid1", "test1")
		assert.NoError(t, err)

		err = store.Add(&core.NewFile{FileId: "2", GroupId: "g2", SyncVersion: 2, EncryptMeta: "B1F2G3", Name: "Budget2"})
		assert.NoError(t, err)
		err = store.UpdateEncryption("2", "salt2", "keyid2", "test2")
		assert.NoError(t, err)

		err = store.Add(&core.NewFile{FileId: "3", GroupId: "g3", SyncVersion: 3, EncryptMeta: "B4F7G9", Name: "Budget3"})
		assert.NoError(t, err)
		err = store.UpdateEncryption("3", "salt3", "keyid3", "test3")
		assert.NoError(t, err)

		files, err := store.All()

		assert.NoError(t, err)
		assert.Equal(t, 3, len(files))
		assert.Equal(t, &core.File{FileId: "1", GroupId: "g1", SyncVersion: 1, EncryptMeta: "A1B2C3", EncryptSalt: "salt1", EncryptKeyId: "keyid1", EncryptTest: "test1", Deleted: false, Name: "Budget1"}, files[0])
		assert.Equal(t, &core.File{FileId: "2", GroupId: "g2", SyncVersion: 2, EncryptMeta: "B1F2G3", EncryptSalt: "salt2", EncryptKeyId: "keyid2", EncryptTest: "test2", Deleted: false, Name: "Budget2"}, files[1])
		assert.Equal(t, &core.File{FileId: "3", GroupId: "g3", SyncVersion: 3, EncryptMeta: "B4F7G9", EncryptSalt: "salt3", EncryptKeyId: "keyid3", EncryptTest: "test3", Deleted: false, Name: "Budget3"}, files[2])
	})
}

func TestFileStore_Update(t *testing.T) {
	t.Run("given no row with matching id", func(t *testing.T) {
		store, conn := newTestFileStore(t)
		defer conn.Close()

		err := store.Update(&core.File{FileId: "1", GroupId: "g1", SyncVersion: 1, EncryptMeta: "A1B2C3", EncryptSalt: "salt1", EncryptKeyId: "keyid1", EncryptTest: "test1", Deleted: false, Name: "Budget1"})

		assert.ErrorIs(t, err, storage.ErrorNoRecordUpdated)
	})

	t.Run("given row with matching id", func(t *testing.T) {
		store, conn := newTestFileStore(t)
		defer conn.Close()

		err := store.Add(&core.NewFile{FileId: "1", GroupId: "g1", SyncVersion: 1, EncryptMeta: "A1B2C3", Name: "Budget1"})
		assert.NoError(t, err)
		err = store.UpdateEncryption("1", "salt1", "keyid1", "test1")
		assert.NoError(t, err)

		err = store.Update(&core.File{FileId: "1", GroupId: "g1", SyncVersion: 2, EncryptMeta: "X9Y6Z7", EncryptSalt: "salt1", EncryptKeyId: "keyid1", EncryptTest: "test1", Deleted: false, Name: "Budget1"})
		assert.NoError(t, err)

		f, err := store.ForId("1", false)

		assert.NoError(t, err)
		assert.Equal(t, &core.File{FileId: "1", GroupId: "g1", SyncVersion: 2, EncryptMeta: "X9Y6Z7", EncryptSalt: "salt1", EncryptKeyId: "keyid1", EncryptTest: "test1", Deleted: false, Name: "Budget1"}, f)
	})
}

func TestFileStore_Add(t *testing.T) {
	t.Run("add new row", func(t *testing.T) {
		store, conn := newTestFileStore(t)
		defer conn.Close()

		err := store.Add(&core.NewFile{FileId: "1", GroupId: "g1", SyncVersion: 1, EncryptMeta: "A1B2C3", Name: "Budget1"})
		assert.NoError(t, err)
		err = store.UpdateEncryption("1", "salt1", "keyid1", "test1")
		assert.NoError(t, err)

		f, err := store.ForId("1", false)

		assert.NoError(t, err)
		assert.Equal(t, &core.File{FileId: "1", GroupId: "g1", SyncVersion: 1, EncryptMeta: "A1B2C3", EncryptSalt: "salt1", EncryptKeyId: "keyid1", EncryptTest: "test1", Deleted: false, Name: "Budget1"}, f)
	})
}

func TestFileStore_ClearGroup(t *testing.T) {
	t.Run("given no row with matching id", func(t *testing.T) {
		store, conn := newTestFileStore(t)
		defer conn.Close()

		err := store.ClearGroup("1")

		assert.ErrorIs(t, err, storage.ErrorNoRecordUpdated)
	})

	t.Run("given row with matching id", func(t *testing.T) {
		store, conn := newTestFileStore(t)
		defer conn.Close()

		err := store.Add(&core.NewFile{FileId: "1", GroupId: "g1", SyncVersion: 1, EncryptMeta: "A1B2C3", Name: "Budget1"})
		assert.NoError(t, err)
		err = store.UpdateEncryption("1", "salt1", "keyid1", "test1")
		assert.NoError(t, err)

		err = store.ClearGroup("1")
		assert.NoError(t, err)

		// TODO: The snuppet below can be used when File.GroupId type is made NULLABLE
		// f, err := store.ForId("1", false)

		// assert.NoError(t, err)
		// assert.Equal(t, &core.File{FileId: "1", GroupId: nil, SyncVersion: 1, EncryptMeta: "A1B2C3", EncryptSalt: "salt1", EncryptKeyId: "keyid1", EncryptTest: "test1", Deleted: false, Name: "Budget1"}, f)
	})
}

func TestFileStore_Delete(t *testing.T) {
	t.Run("given no row with matching id", func(t *testing.T) {
		store, conn := newTestFileStore(t)
		defer conn.Close()

		err := store.Delete("1")

		assert.ErrorIs(t, err, storage.ErrorNoRecordUpdated)
	})

	t.Run("given row with matching id", func(t *testing.T) {
		store, conn := newTestFileStore(t)
		defer conn.Close()

		err := store.Add(&core.NewFile{FileId: "1", GroupId: "g1", SyncVersion: 1, EncryptMeta: "A1B2C3", Name: "Budget1"})
		assert.NoError(t, err)
		err = store.UpdateEncryption("1", "salt1", "keyid1", "test1")
		assert.NoError(t, err)

		err = store.Delete("1")
		assert.NoError(t, err)

		f, err := store.ForId("1", true)

		assert.NoError(t, err)
		assert.Equal(t, &core.File{FileId: "1", GroupId: "g1", SyncVersion: 1, EncryptMeta: "A1B2C3", EncryptSalt: "salt1", EncryptKeyId: "keyid1", EncryptTest: "test1", Deleted: true, Name: "Budget1"}, f)
	})
}

func TestFileStore_UpdateName(t *testing.T) {
	t.Run("given no row with matching id", func(t *testing.T) {
		store, conn := newTestFileStore(t)
		defer conn.Close()

		err := store.UpdateName("1", "My budget")

		assert.ErrorIs(t, err, storage.ErrorNoRecordUpdated)
	})

	t.Run("given row with matching id", func(t *testing.T) {
		store, conn := newTestFileStore(t)
		defer conn.Close()

		err := store.Add(&core.NewFile{FileId: "1", GroupId: "g1", SyncVersion: 1, EncryptMeta: "A1B2C3", Name: "Budget1"})
		assert.NoError(t, err)
		err = store.UpdateEncryption("1", "salt1", "keyid1", "test1")
		assert.NoError(t, err)

		err = store.UpdateName("1", "My budget")
		assert.NoError(t, err)

		f, err := store.ForId("1", false)

		assert.NoError(t, err)
		assert.Equal(t, &core.File{FileId: "1", GroupId: "g1", SyncVersion: 1, EncryptMeta: "A1B2C3", EncryptSalt: "salt1", EncryptKeyId: "keyid1", EncryptTest: "test1", Deleted: false, Name: "My budget"}, f)
	})
}

func TestFileStore_UpdateGroup(t *testing.T) {
	t.Run("given no row with matching id", func(t *testing.T) {
		store, conn := newTestFileStore(t)
		defer conn.Close()

		err := store.UpdateGroup("1", "gnew")

		assert.ErrorIs(t, err, storage.ErrorNoRecordUpdated)
	})

	t.Run("given row with matching id", func(t *testing.T) {
		store, conn := newTestFileStore(t)
		defer conn.Close()

		err := store.Add(&core.NewFile{FileId: "1", GroupId: "g1", SyncVersion: 1, EncryptMeta: "A1B2C3", Name: "Budget1"})
		assert.NoError(t, err)
		err = store.UpdateEncryption("1", "salt1", "keyid1", "test1")
		assert.NoError(t, err)

		err = store.UpdateGroup("1", "gnew")
		assert.NoError(t, err)

		f, err := store.ForId("1", false)

		assert.NoError(t, err)
		assert.Equal(t, &core.File{FileId: "1", GroupId: "gnew", SyncVersion: 1, EncryptMeta: "A1B2C3", EncryptSalt: "salt1", EncryptKeyId: "keyid1", EncryptTest: "test1", Deleted: false, Name: "Budget1"}, f)
	})
}

func TestFileStore_UpdateEncryption(t *testing.T) {
	t.Run("given no row with matching id", func(t *testing.T) {
		store, conn := newTestFileStore(t)
		defer conn.Close()

		err := store.UpdateEncryption("1", "saltNew", "keyidNew", "testNew")

		assert.ErrorIs(t, err, storage.ErrorNoRecordUpdated)
	})

	t.Run("given row with matching id", func(t *testing.T) {
		store, conn := newTestFileStore(t)
		defer conn.Close()

		err := store.Add(&core.NewFile{FileId: "1", GroupId: "g1", SyncVersion: 1, EncryptMeta: "A1B2C3", Name: "Budget1"})
		assert.NoError(t, err)
		err = store.UpdateEncryption("1", "salt1", "keyid1", "test1")
		assert.NoError(t, err)

		err = store.UpdateEncryption("1", "saltNew", "keyidNew", "testNew")
		assert.NoError(t, err)

		f, err := store.ForId("1", false)

		assert.NoError(t, err)
		assert.Equal(t, &core.File{FileId: "1", GroupId: "g1", SyncVersion: 1, EncryptMeta: "A1B2C3", EncryptSalt: "saltNew", EncryptKeyId: "keyidNew", EncryptTest: "testNew", Deleted: false, Name: "Budget1"}, f)
	})
}
