package sqlite

import (
	"database/sql"
	"errors"

	"github.com/nathanjisaac/actual-server-go/internal/core"
	internal_errors "github.com/nathanjisaac/actual-server-go/internal/errors"
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
	_, _, err := ms.connection.Mutate(
		"INSERT INTO messages_merkles (id, merkle) VALUES (?, ?) ON CONFLICT (id) DO UPDATE SET merkle = ?",
		message.MerkleID,
		message.Merkle,
		message.Merkle,
	)
	if err != nil {
		return err
	}
	return nil
}

func (ms *MerkleStore) GetForGroup(groupID string) (*core.MerkleMessage, error) {
	row, err := ms.connection.First("SELECT * FROM messages_merkles WHERE id = ?", groupID)
	if err != nil {
		return nil, err
	}

	var msg core.MerkleMessage
	if err = row.Scan(&msg.MerkleID, &msg.Merkle); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, internal_errors.ErrStorageRecordNotFound
		}
		return nil, err
	}

	return &msg, nil
}
