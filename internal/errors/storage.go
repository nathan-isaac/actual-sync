package errors

import "errors"

var (
	StorageErrorRecordNotFound  = errors.New("record not found")
	StorageErrorNoRecordUpdated = errors.New("no record updated")
)
