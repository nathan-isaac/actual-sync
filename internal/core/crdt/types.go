package crdt

type Merkle interface {
	ToJSONString() (string, error)
	Insert(Timestamp)
	Prune() Merkle
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
