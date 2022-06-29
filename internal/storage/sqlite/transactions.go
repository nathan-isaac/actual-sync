package sqlite

import (
	"database/sql"
	"encoding/json"

	"github.com/nathanjisaac/actual-server-go/internal/core"
	"github.com/nathanjisaac/actual-server-go/internal/core/crdt"
	"github.com/nathanjisaac/actual-server-go/internal/core/crdt/merkle"
	"github.com/nathanjisaac/actual-server-go/internal/core/crdt/timestamp"
	"github.com/nathanjisaac/actual-server-go/internal/routes/syncpb"
)

func AddNewMessagesTransaction(db *Connection, messages []*syncpb.MessageEnvelope) (crdt.Merkle, error) {
	merkleTrie := merkle.NewMerkle(0)
	err := db.Transaction(func(tx *sql.Tx) error {
		trie, err := getMerkle(tx)
		if err != nil {
			return err
		}

		if len(messages) > 0 {
			for _, msg := range messages {
				err := updateBinaryMerkleStore(tx, msg, trie)
				if err != nil {
					return err
				}
			}
		}

		prunedTrie := trie.Prune().(*merkle.Merkle)

		err = updateMessagesStore(tx, prunedTrie)
		if err != nil {
			return err
		}
		merkleTrie = prunedTrie
		return nil
	})
	if err != nil {
		return nil, err
	}

	return merkleTrie, nil
}

func getMerkle(tx *sql.Tx) (*merkle.Merkle, error) {
	stmt, err := tx.Prepare("SELECT * FROM messages_merkles")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRow()

	var msg core.MerkleMessage
	if err = row.Scan(&msg.MerkleId, &msg.Merkle); err != nil {
		if err == sql.ErrNoRows {
			return merkle.NewMerkle(0), nil
		}
		return nil, err
	}

	var merkleMap map[string]interface{}
	err = json.Unmarshal([]byte(msg.Merkle), &merkleMap)
	if err != nil {
		return nil, err
	}

	merkle := merkle.NewMerkleFromMap(merkleMap)

	return merkle, nil
}

func updateBinaryMerkleStore(tx *sql.Tx, msg *syncpb.MessageEnvelope, trie crdt.Merkle) error {
	stmt, err := tx.Prepare("INSERT OR IGNORE INTO messages_binary (timestamp, is_encrypted, content) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}

	defer stmt.Close()

	result, err := stmt.Exec(msg.Timestamp, msg.IsEncrypted, msg.Content)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows > 0 {
		ts, err := timestamp.ParseTimestamp(msg.Timestamp)
		if err != nil {
			return err
		}
		trie.Insert(ts)
		return nil
	}

	return nil
}

func updateMessagesStore(tx *sql.Tx, trie *merkle.Merkle) error {
	stmt, err := tx.Prepare("INSERT INTO messages_merkles (id, merkle) VALUES (1, ?) ON CONFLICT (id) DO UPDATE SET merkle = ?")
	if err != nil {
		return err
	}

	defer stmt.Close()

	trieString, err := trie.ToJSONString()
	if err != nil {
		return err
	}
	_, err = stmt.Exec(trieString, trieString)
	if err != nil {
		return err
	}

	return nil
}
