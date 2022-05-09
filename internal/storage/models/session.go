package models

import "database/sql"

type Session struct {
	Token string
}

func (s *Session) FromRow(row *sql.Row) {
	row.Scan(&s.Token)
}

func (s *Session) FromRows(rows *sql.Row) {
	rows.Scan(&s.Token)
}
