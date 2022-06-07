package merkle

import (
	"sort"
	"strconv"

	"github.com/nathanjisaac/actual-server-go/internal/core/crdt"
)

type Merkle struct {
	Hash     uint32
	Children map[string]*Merkle
}

func NewMerkle(hash uint32) *Merkle {
	return &Merkle{Hash: hash, Children: map[string]*Merkle{}}
}

func deepCopyMerkle(merkle *Merkle) *Merkle {
	if len(merkle.Children) == 0 {
		return &Merkle{Hash: merkle.Hash, Children: map[string]*Merkle{}}
	}

	children := map[string]*Merkle{}
	for k := range merkle.Children {
		children[k] = deepCopyMerkle(merkle.Children[k])
	}

	return &Merkle{Hash: merkle.Hash, Children: children}
}

func (trie *Merkle) GetKeys() []string {
	j := 0
	keys := make([]string, len(trie.Children))
	for k := range trie.Children {
		keys[j] = k
		j++
	}
	return keys
}

func (trie *Merkle) insertKey(key string, hash uint32) *Merkle {
	newTrie := deepCopyMerkle(trie)
	if len(key) == 0 {
		return newTrie
	}
	c := string(key[0])
	n := newTrie.Children[c]
	newTrie.Children[c] = newTrie.insertKey(key[1:], hash)
	if n != nil {
		newTrie.Children[c].Hash = n.Hash ^ hash
	} else {
		newTrie.Children[c].Hash = hash
	}
	return newTrie
}

func (trie *Merkle) Insert(timestamp crdt.Timestamp) *Merkle {
	hash := timestamp.Hash()
	key := strconv.FormatInt((timestamp.GetMillis() / (1000 * 60)), 3)

	newTrie := deepCopyMerkle(trie)
	newTrie.Hash = trie.Hash ^ hash
	return newTrie.insertKey(key, hash)
}

func (trie *Merkle) Prune() *Merkle {
	n := 2

	// Checking if empty
	if len(trie.Children) == 0 {
		return NewMerkle(trie.Hash)
	}

	keys := trie.GetKeys()
	sort.Strings(keys)

	newTrie := NewMerkle(trie.Hash)
	sliceRange := 0
	if (len(keys) - n) > 0 {
		sliceRange = len(keys) - n
	}
	for k := range keys[sliceRange:] {
		newTrie.Children[keys[k]] = trie.Children[keys[k]].Prune()
	}
	return newTrie
}
