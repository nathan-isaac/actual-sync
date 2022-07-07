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
