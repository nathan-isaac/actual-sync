package storage_sqlite

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

type WrappedDatabase struct {
	DB *sql.DB
}

func (wdb *WrappedDatabase) All(sql string, params ...any) *sql.Rows {
	stmt, _ := wdb.DB.Prepare(sql)
	defer stmt.Close()
	rows, _ := stmt.Query(params...)
	defer rows.Close()
	return rows
}

func (wdb *WrappedDatabase) First(sql string, params ...any) *sql.Row {
	stmt, _ := wdb.DB.Prepare(sql)
	defer stmt.Close()
	row := stmt.QueryRow(params...)
	if row.Err() != nil {
		return nil
	}
	return row
}

func (wdb *WrappedDatabase) Exec(sql string) {
	wdb.DB.Exec(sql)
}

func (wdb *WrappedDatabase) Mutate(sql string, params ...any) (int64, int64) {
	stmt, _ := wdb.DB.Prepare(sql)
	defer stmt.Close()
	result, _ := stmt.Exec(params...)
	rows, _ := result.RowsAffected()
	lastId, _ := result.LastInsertId()
	return rows, lastId
}

func (wdb *WrappedDatabase) Close() {
	wdb.DB.Close()
}

func OpenDatabase(filename string) WrappedDatabase {
	db, err := sql.Open("sqlite", filename)
	if err != nil {
		log.Fatal(err)
	}
	return WrappedDatabase{DB: db}
}
