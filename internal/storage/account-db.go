package storage

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/nathanjisaac/actual-server-go/internal/core"
	storage_sqlite "github.com/nathanjisaac/actual-server-go/internal/storage/sqlite"
)

var accountDb storage_sqlite.WrappedDatabase

func GetAccountDB(config core.Config) storage_sqlite.WrappedDatabase {
	if accountDb.DB == nil {
		dbPath := filepath.Join(config.ServerFiles, "account.sqlite")
		_, err := os.Stat(dbPath)
		dbExists := !errors.Is(err, os.ErrNotExist)

		accountDb = storage_sqlite.OpenDatabase(dbPath)

		if !dbExists {
			sqlFilepath, _ := filepath.Abs("internal/storage/sql-scripts/account.sql")
			content, _ := ioutil.ReadFile(sqlFilepath)
			initSql := string(content)
			accountDb.Exec(initSql)
		}
	}
	return accountDb
}
