package errors

import "errors"

var (
	ErrTimestampUnableToParse = errors.New("unable to parse timestamp")
)
