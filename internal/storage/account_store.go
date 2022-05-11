package storage

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/nathanjisaac/actual-server-go/internal/core"
	"github.com/nathanjisaac/actual-server-go/internal/storage/sqlite"
)

var accountStore *sqlite.Connection

func GetAccountDB(config core.Config) *sqlite.Connection {
	if accountStore == nil {
		dbPath := filepath.Join(config.ServerFiles, "account.sqlite")
		_, err := os.Stat(dbPath)
		dbExists := !errors.Is(err, os.ErrNotExist)

		accountStore, err = sqlite.NewConnection(dbPath)
		if err != nil {
			log.Fatal(err)
		}

		if !dbExists {
			sqlFilepath, _ := filepath.Abs("internal/storage/sqlite/migrations/init_account.sql")
			content, _ := ioutil.ReadFile(sqlFilepath)
			initSql := string(content)
			accountStore.Exec(initSql)
		}
	}

	return accountStore
}
