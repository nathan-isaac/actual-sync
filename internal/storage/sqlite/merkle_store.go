package sqlite

import (
	"database/sql"

	"github.com/nathanjisaac/actual-server-go/internal/core"
	"github.com/nathanjisaac/actual-server-go/internal/storage"
)

type MerkleStore struct {
	connection *Connection
}

func NewMerkleStore(connection *Connection) *MerkleStore {
	return &MerkleStore{
		connection: connection,
	}
}

func (ms *MerkleStore) Add(message core.MerkleMessage) error {
	_, _, err := ms.connection.Mutate("INSERT INTO messages_merkles (id, merkle) VALUES (?, ?) ON CONFLICT (id) DO UPDATE SET merkle = ?", message.MerkleId, message.Merkle, message.Merkle)
	if err != nil {
		return err
	}
	return nil
}

func (ms *MerkleStore) GetForGroup(groupId string) (*core.MerkleMessage, error) {
	row, err := ms.connection.First("SELECT * FROM messages_merkles WHERE id = ?", groupId)
	if err != nil {
		return nil, err
	}

	var msg core.MerkleMessage
	if err = row.Scan(&msg.MerkleId, &msg.Merkle); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.ErrorRecordNotFound
		}
		return nil, err
	}

	return &msg, nil
}
