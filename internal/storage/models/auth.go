package models

import "database/sql"

type Auth struct {
	Password string
}

func (a *Auth) FromRow(row *sql.Row) {
	row.Scan(&a.Password)
}

func (a *Auth) FromRows(rows *sql.Rows) {
	rows.Scan(&a.Password)
}
