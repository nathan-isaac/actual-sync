package merkle_test

import (
	"testing"

	"github.com/nathanjisaac/actual-server-go/internal/core/crdt/merkle"
	"github.com/nathanjisaac/actual-server-go/internal/core/crdt/timestamp"
	"github.com/stretchr/testify/assert"
)

func TestMerkle_GetKeys(t *testing.T) {
	t.Run("return one key", func(t *testing.T) {
		trie := merkle.NewMerkle(2)
		var time *timestamp.Timestamp
		time, _ = timestamp.ParseTimestamp("2018-11-12T13:21:40.122Z-0000-0123456789ABCDEF")
		trie = trie.Insert(time)
		assert.Equal(t, []string{"1"}, trie.GetKeys())
	})
}
