package errors

import "errors"

var (
	ErrStorageRecordNotFound  = errors.New("record not found")
	ErrStorageNoRecordUpdated = errors.New("no record updated")
)
