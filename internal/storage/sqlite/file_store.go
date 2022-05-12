package sqlite

import (
	"database/sql"

	"github.com/nathanjisaac/actual-server-go/internal/core"
	"github.com/nathanjisaac/actual-server-go/internal/storage"
)

type FileStore struct {
	connection *Connection
}

func NewFileStore(connection *Connection) *FileStore {
	return &FileStore{
		connection: connection,
	}
}

func (fs *FileStore) ForId(id core.FileId, deleted bool) (*core.File, error) {
	row, err := fs.connection.First("SELECT * FROM files WHERE id = ? AND deleted = ?", id, deleted)
	if err != nil {
		return nil, err
	}

	var f core.File

	if err = row.Scan(&f.FileId, &f.GroupId, &f.SyncVersion, &f.EncryptMeta, &f.EncryptKeyId, &f.EncryptSalt, &f.EncryptTest, &f.Deleted, &f.Name); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.ErrorRecordNotFound
		}
		return nil, err
	}

	return &f, nil
}

func (fs *FileStore) All() ([]*core.File, error) {
	rows, err := fs.connection.All("SELECT * FROM files")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	files := make([]*core.File, 0)
	for rows.Next() {
		var f core.File
		if err := rows.Scan(&f.FileId, &f.GroupId, &f.SyncVersion, &f.EncryptMeta, &f.EncryptKeyId, &f.EncryptSalt, &f.EncryptTest, &f.Deleted, &f.Name); err != nil {
			return nil, err
		}
		files = append(files, &f)
	}

	return files, nil
}

func (fs *FileStore) Update(file *core.File) error {
	rows, _, err := fs.connection.Mutate("UPDATE files SET sync_version = ?, encrypt_meta = ? WHERE id = ?", file.SyncVersion, file.EncryptMeta, file.FileId)
	if err != nil {
		return err
	} else if rows == 0 {
		return storage.ErrorNoRecordUpdated
	}

	return nil
}

func (fs *FileStore) Add(file *core.NewFile) error {
	_, _, err := fs.connection.Mutate("INSERT INTO files (id, group_id, sync_version, name, encrypt_meta) VALUES (?, ?, ?, ?, ?)", file.FileId, file.GroupId, file.SyncVersion, file.Name, file.EncryptMeta)
	if err != nil {
		return err
	}

	return nil
}

func (fs *FileStore) ClearGroup(id core.FileId) error {
	rows, _, err := fs.connection.Mutate("UPDATE files SET group_id = NULL WHERE id = ?", id)
	if err != nil {
		return err
	} else if rows == 0 {
		return storage.ErrorNoRecordUpdated
	}

	return nil
}

func (fs *FileStore) Delete(id core.FileId) error {
	rows, _, err := fs.connection.Mutate("UPDATE files SET deleted = TRUE WHERE id = ?", id)
	if err != nil {
		return err
	} else if rows == 0 {
		return storage.ErrorNoRecordUpdated
	}

	return nil
}

func (fs *FileStore) UpdateName(id core.FileId, name string) error {
	rows, _, err := fs.connection.Mutate("UPDATE files SET name = ? WHERE id = ?", name, id)
	if err != nil {
		return err
	} else if rows == 0 {
		return storage.ErrorNoRecordUpdated
	}

	return nil
}

func (fs *FileStore) UpdateGroup(id core.FileId, groupId string) error {
	rows, _, err := fs.connection.Mutate("UPDATE files SET group_id = ? WHERE id = ?", groupId, id)
	if err != nil {
		return err
	} else if rows == 0 {
		return storage.ErrorNoRecordUpdated
	}

	return nil
}

func (fs *FileStore) UpdateEncryption(id core.FileId, salt, keyId, test string) error {
	rows, _, err := fs.connection.Mutate("UPDATE files SET encrypt_salt = ?, encrypt_keyid = ?, encrypt_test = ? WHERE id = ?", salt, keyId, test, id)
	if err != nil {
		return err
	} else if rows == 0 {
		return storage.ErrorNoRecordUpdated
	}

	return nil
}
