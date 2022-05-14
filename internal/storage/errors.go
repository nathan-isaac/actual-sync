package storage

import "errors"

var (
	ErrorRecordNotFound  = errors.New("record not found")
	ErrorNoRecordUpdated = errors.New("no record updated")
)
