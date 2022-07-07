package sqlite

import (
	"database/sql"
	"errors"

	"github.com/nathanjisaac/actual-server-go/internal/core"
	internal_errors "github.com/nathanjisaac/actual-server-go/internal/errors"
)

type FileStore struct {
	connection *Connection
}

func NewFileStore(connection *Connection) *FileStore {
	return &FileStore{
		connection: connection,
	}
}

func (fs *FileStore) Count() (int, error) {
	row, err := fs.connection.First("SELECT count(*) FROM files")

	if err != nil {
		return 0, err
	}

	var count int

	if err = row.Scan(&count); err != nil {
		return 0, internal_errors.ErrStorageRecordNotFound
	}

	return count, nil
}

func (fs *FileStore) ForID(id core.FileID) (*core.File, error) {
	row, err := fs.connection.First("SELECT * FROM files WHERE id = ?", id)
	if err != nil {
		return nil, err
	}

	var f core.File
	var gid sql.NullString
	var encryptKey sql.NullString
	var encryptSalt sql.NullString
	var encryptTest sql.NullString

	if err = row.Scan(
		&f.FileID,
		&gid,
		&f.SyncVersion,
		&f.EncryptMeta,
		&encryptKey,
		&encryptSalt,
		&encryptTest,
		&f.Deleted,
		&f.Name,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, internal_errors.ErrStorageRecordNotFound
		}
		return nil, err
	}
	if gid.Valid {
		f.GroupID = gid.String
	}
	if encryptKey.Valid {
		f.EncryptKeyID = encryptKey.String
	}
	if encryptSalt.Valid {
		f.EncryptSalt = encryptSalt.String
	}
	if encryptTest.Valid {
		f.EncryptTest = encryptTest.String
	}

	return &f, nil
}

func (fs *FileStore) ForIDAndDelete(id core.FileID, deleted bool) (*core.File, error) {
	row, err := fs.connection.First("SELECT * FROM files WHERE id = ? AND deleted = ?", id, deleted)
	if err != nil {
		return nil, err
	}

	var f core.File
	var gid sql.NullString
	var encryptKey sql.NullString
	var encryptSalt sql.NullString
	var encryptTest sql.NullString

	if err = row.Scan(
		&f.FileID,
		&gid,
		&f.SyncVersion,
		&f.EncryptMeta,
		&encryptKey,
		&encryptSalt,
		&encryptTest,
		&f.Deleted,
		&f.Name,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, internal_errors.ErrStorageRecordNotFound
		}
		return nil, err
	}
	if gid.Valid {
		f.GroupID = gid.String
	}
	if encryptKey.Valid {
		f.EncryptKeyID = encryptKey.String
	}
	if encryptSalt.Valid {
		f.EncryptSalt = encryptSalt.String
	}
	if encryptTest.Valid {
		f.EncryptTest = encryptTest.String
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
		var gid sql.NullString
		var encryptKey sql.NullString
		var encryptSalt sql.NullString
		var encryptTest sql.NullString

		if err := rows.Scan(
			&f.FileID,
			&gid,
			&f.SyncVersion,
			&f.EncryptMeta,
			&encryptKey,
			&encryptSalt,
			&encryptTest,
			&f.Deleted,
			&f.Name,
		); err != nil {
			return nil, err
		}
		if gid.Valid {
			f.GroupID = gid.String
		}
		if encryptKey.Valid {
			f.EncryptKeyID = encryptKey.String
		}
		if encryptSalt.Valid {
			f.EncryptSalt = encryptSalt.String
		}
		if encryptTest.Valid {
			f.EncryptTest = encryptTest.String
		}

		files = append(files, &f)
	}

	return files, nil
}

func (fs *FileStore) Update(fileID string, syncVersion int16, encryptMeta string, name string) error {
	rows, _, err := fs.connection.Mutate(
		"UPDATE files SET sync_version = ?, encrypt_meta = ?, name = ? WHERE id = ?",
		syncVersion,
		encryptMeta,
		name,
		fileID,
	)
	if err != nil {
		return err
	} else if rows == 0 {
		return internal_errors.ErrStorageNoRecordUpdated
	}

	return nil
}

func (fs *FileStore) Add(file *core.NewFile) error {
	_, _, err := fs.connection.Mutate(
		"INSERT INTO files (id, group_id, sync_version, name, encrypt_meta) VALUES (?, ?, ?, ?, ?)",
		file.FileID,
		file.GroupID,
		file.SyncVersion,
		file.Name,
		file.EncryptMeta,
	)
	if err != nil {
		return err
	}

	return nil
}

func (fs *FileStore) ClearGroup(id core.FileID) error {
	rows, _, err := fs.connection.Mutate("UPDATE files SET group_id = NULL WHERE id = ?", id)
	if err != nil {
		return err
	} else if rows == 0 {
		return internal_errors.ErrStorageNoRecordUpdated
	}

	return nil
}

func (fs *FileStore) Delete(id core.FileID) error {
	rows, _, err := fs.connection.Mutate("UPDATE files SET deleted = TRUE WHERE id = ?", id)
	if err != nil {
		return err
	} else if rows == 0 {
		return internal_errors.ErrStorageNoRecordUpdated
	}

	return nil
}

func (fs *FileStore) UpdateName(id core.FileID, name string) error {
	rows, _, err := fs.connection.Mutate("UPDATE files SET name = ? WHERE id = ?", name, id)
	if err != nil {
		return err
	} else if rows == 0 {
		return internal_errors.ErrStorageNoRecordUpdated
	}

	return nil
}

func (fs *FileStore) UpdateGroup(id core.FileID, groupID string) error {
	rows, _, err := fs.connection.Mutate("UPDATE files SET group_id = ? WHERE id = ?", groupID, id)
	if err != nil {
		return err
	} else if rows == 0 {
		return internal_errors.ErrStorageNoRecordUpdated
	}

	return nil
}

func (fs *FileStore) UpdateEncryption(id core.FileID, salt, keyID, test string) error {
	rows, _, err := fs.connection.Mutate(
		"UPDATE files SET encrypt_salt = ?, encrypt_keyid = ?, encrypt_test = ? WHERE id = ?",
		salt,
		keyID,
		test,
		id,
	)
	if err != nil {
		return err
	} else if rows == 0 {
		return internal_errors.ErrStorageNoRecordUpdated
	}

	return nil
}
