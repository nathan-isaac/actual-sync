package core

type MerkleMessage struct {
	MerkleID string
	Merkle   string
}

type MerkleStore interface {
	Add(message MerkleMessage) error
	GetForGroup(groupID string) (*MerkleMessage, error)
}
