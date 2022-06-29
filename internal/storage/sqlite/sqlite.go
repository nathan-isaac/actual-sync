package sqlite

import "github.com/nathanjisaac/actual-server-go/internal/core"

type StorageConfig struct {
	ServerData string
	UserData   string
}

func NewAccountStores(dataSource string) (core.Connection, core.PasswordStore, core.TokenStore, core.FileStore, error) {
	db, err := NewAccountConnection(dataSource)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	pwdStore := NewPasswordStore(db)
	tokenStore := NewTokenStore(db)
	fileStore := NewFileStore(db)
	return db, pwdStore, tokenStore, fileStore, nil
}

func NewGroupStores(dataSource string) (core.Connection, core.MerkleStore, core.MessageStore, error) {
	db, err := NewMessageConnection(dataSource)
	if err != nil {
		return nil, nil, nil, err
	}

	merkleDb := NewMerkleStore(db)
	messageDb := NewMessageStore(db)
	return db, merkleDb, messageDb, nil
}
