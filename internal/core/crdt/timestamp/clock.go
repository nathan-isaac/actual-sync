package timestamp

import (
	"github.com/nathanjisaac/actual-server-go/internal/core/crdt"
)

type Clock struct {
	Timestamp MutableTimestamp
	Merkle    crdt.Merkle
}

func NewClock(timestamp MutableTimestamp, merkle crdt.Merkle) *Clock {
	return &Clock{Timestamp: timestamp, Merkle: merkle}
}
