package sqlite

import (
	"github.com/nathanjisaac/actual-server-go/internal/core"
)

type MessageStore struct {
	connection *Connection
}

func NewMessageStore(connection *Connection) *MessageStore {
	return &MessageStore{
		connection: connection,
	}
}

func (ms *MessageStore) Add(message core.BinaryMessage) (bool, error) {
	rowsAffected, _, err := ms.connection.Mutate(
		"INSERT OR IGNORE INTO messages_binary (timestamp, is_encrypted, content) VALUES (?, ?, ?)",
		message.Timestamp,
		message.IsEncrypted,
		message.Content,
	)
	if err != nil {
		return false, err
	}
	if rowsAffected > 0 {
		return true, nil
	}
	return false, nil
}

func (ms *MessageStore) GetSince(timestamp string) ([]*core.BinaryMessage, error) {
	rows, err := ms.connection.All("SELECT * FROM messages_binary WHERE timestamp > ? ORDER BY timestamp", timestamp)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := make([]*core.BinaryMessage, 0)
	for rows.Next() {
		var msg core.BinaryMessage

		if err := rows.Scan(&msg.Timestamp, &msg.IsEncrypted, &msg.Content); err != nil {
			return nil, err
		}

		messages = append(messages, &msg)
	}

	return messages, nil
}
