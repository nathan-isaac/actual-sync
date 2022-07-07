package storage

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nathanjisaac/actual-server-go/internal/core"
	"github.com/nathanjisaac/actual-server-go/internal/core/crdt"
	"github.com/nathanjisaac/actual-server-go/internal/routes/syncpb"
	"github.com/nathanjisaac/actual-server-go/internal/storage/sqlite"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

type Options struct {
	DataPath       string
	ServerDataPath string
	UserDataPath   string
}

func GenerateStorageConfig(storage string, options Options) core.StorageConfig {
	switch storage {
	case string(core.Sqlite):
		if options.ServerDataPath == "" {
			options.ServerDataPath = filepath.Join(options.DataPath, "server-files")
		} else {
			if !filepath.IsAbs(options.ServerDataPath) {
				path, err := filepath.Abs(options.ServerDataPath)
				cobra.CheckErr(err)
				options.ServerDataPath = path
			}
		}
		if options.UserDataPath == "" {
			options.UserDataPath = filepath.Join(options.DataPath, "user-files")
		} else {
			if !filepath.IsAbs(options.UserDataPath) {
				path, err := filepath.Abs(options.UserDataPath)
				cobra.CheckErr(err)
				options.UserDataPath = path
			}
		}

		fs := afero.NewOsFs()

		err := fs.MkdirAll(options.ServerDataPath, os.ModePerm)
		cobra.CheckErr(err)

		err = fs.MkdirAll(options.UserDataPath, os.ModePerm)
		cobra.CheckErr(err)
		return sqlite.StorageConfig{
			ServerData: options.ServerDataPath,
			UserData:   options.UserDataPath,
		}
	default:
		cobra.CheckErr("Invalid storage type!")
	}
	return nil
}

func NewAccountStores(storageType core.StorageType, config core.StorageConfig) (core.Connection, core.PasswordStore, core.TokenStore, core.FileStore, error) {
	switch storageType {
	case core.Sqlite:
		return sqlite.NewAccountStores(filepath.Join(config.(sqlite.StorageConfig).ServerData, "account.sqlite"))
	default:
		// Default is set to Sqlite
		return sqlite.NewAccountStores(filepath.Join(config.(sqlite.StorageConfig).ServerData, "account.sqlite"))
	}
}

func NewGroupStores(storageType core.StorageType, config core.StorageConfig, fileId core.FileId) (core.Connection, core.MerkleStore, core.MessageStore, error) {
	fileName := fmt.Sprintf("%s.sqlite", fileId)
	switch storageType {
	case core.Sqlite:
		return sqlite.NewGroupStores(filepath.Join(config.(sqlite.StorageConfig).UserData, fileName))
	default:
		// Default is set to Sqlite
		return sqlite.NewGroupStores(filepath.Join(config.(sqlite.StorageConfig).UserData, fileName))
	}
}

func AddNewMessagesTransaction(storageType core.StorageType, db core.Connection, messages []*syncpb.MessageEnvelope) (crdt.Merkle, error) {
	switch storageType {
	case core.Sqlite:
		return sqlite.AddNewMessagesTransaction(db.(*sqlite.Connection), messages)
	default:
		// Default is set to Sqlite
		return sqlite.AddNewMessagesTransaction(db.(*sqlite.Connection), messages)
	}
}
