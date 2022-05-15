package crdt

type Merkle interface {
	GetKeys(Merkle) []string
	Insert(Timestamp) *Merkle
	Prune(Timestamp) *Merkle
}

type Timestamp interface {
	ToString() string
	GetMillis() int64
	GetCounter() int64
	GetNode() string
	Hash() uint32
}

type MutableTimestamp interface {
	Timestamp
	SetMillis(int64)
	SetCounter(int64)
	SetNode(string)
}
