package storage

import (
	"github.com/nathanjisaac/actual-server-go/internal/core"
	"github.com/nathanjisaac/actual-server-go/internal/storage/sqlite"
)

func NewSqliteAccountStores(dataSource string) (*sqlite.Connection, core.PasswordStore, core.TokenStore, core.FileStore, error) {
	db, err := sqlite.NewAccountConnection(dataSource)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	pwdStore := sqlite.NewPasswordStore(db)
	tokenStore := sqlite.NewTokenStore(db)
	fileStore := sqlite.NewFileStore(db)
	return db, pwdStore, tokenStore, fileStore, nil
}

func NewSqliteGroupStores(dataSource string) (*sqlite.Connection, core.MerkleStore, core.MessageStore, error) {
	db, err := sqlite.NewMessageConnection(dataSource)
	if err != nil {
		return nil, nil, nil, err
	}

	merkleDb := sqlite.NewMerkleStore(db)
	messageDb := sqlite.NewMessageStore(db)
	return db, merkleDb, messageDb, nil
}
