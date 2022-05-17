package core

type BinaryMessage struct {
	Timestamp   string
	IsEncrypted bool
	Content     []byte
}

type MessageStore interface {
	Add(message BinaryMessage) (bool, error)
	GetSince(timestamp string) ([]*BinaryMessage, error)
}

type MerkleMessage struct {
	MerkleId string
	Merkle   string
}

type MerkleStore interface {
	Add(message MerkleMessage) error
	GetForGroup(groupId string) (*MerkleMessage, error)
}
