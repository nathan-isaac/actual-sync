package models

import "database/sql"

type File struct {
	Id            string
	Group_id      string
	Sync_version  int16
	Encrypt_meta  string
	Encrypt_keyid string
	Encrypt_salt  string
	Encrypt_test  string
	Deleted       bool
	Name          string
}

func (f *File) FromRow(row *sql.Row) {
	row.Scan(&f.Id, &f.Group_id, &f.Sync_version, &f.Encrypt_meta, &f.Encrypt_keyid, &f.Encrypt_salt, &f.Encrypt_test, &f.Deleted, &f.Name)
}

func (f *File) FromRows(rows *sql.Rows) {
	rows.Scan(&f.Id, &f.Group_id, &f.Sync_version, &f.Encrypt_meta, &f.Encrypt_keyid, &f.Encrypt_salt, &f.Encrypt_test, &f.Deleted, &f.Name)
}
