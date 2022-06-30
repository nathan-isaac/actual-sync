package routes

import (
	"github.com/nathanjisaac/actual-server-go/internal/core"
	"github.com/nathanjisaac/actual-server-go/internal/routes/syncpb"
	"github.com/nathanjisaac/actual-server-go/internal/storage"
)

func encryptedSync(
	since string,
	messages []*syncpb.MessageEnvelope,
	fileID core.FileID,
	storageType core.StorageType,
	storageConfig core.StorageConfig,
) (string, []*syncpb.MessageEnvelope, error) {
	db, _, msgStore, err := storage.NewGroupStores(storageType, storageConfig, fileID)
	if err != nil {
		return "", nil, err
	}
	if db != nil {
		defer db.Close()
	}

	newMessages, err := msgStore.GetSince(since)
	if err != nil {
		return "", nil, err
	}
	pbNewMessages := make([]*syncpb.MessageEnvelope, len(newMessages))
	for i, msg := range newMessages {
		pbNewMessages[i] = &syncpb.MessageEnvelope{
			Timestamp:   msg.Timestamp,
			IsEncrypted: msg.IsEncrypted,
			Content:     msg.Content,
		}
	}

	trie, err := storage.AddNewMessagesTransaction(storageType, db, messages)
	if err != nil {
		return "", nil, err
	}

	merkleString, err := trie.ToJSONString()
	if err != nil {
		return "", nil, err
	}

	return merkleString, pbNewMessages, nil
}
