package merkle

import (
	"encoding/json"
	"math"
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

func NewMerkleFromMap(merkle map[string]interface{}) *Merkle {
	if merkle != nil {
		if merkle["hash"] == nil {
			return &Merkle{}
		}
		children := map[string]*Merkle{}
		for key := range merkle {
			if key != "hash" {
				children[key] = NewMerkleFromMap(merkle[key].(map[string]interface{}))
			}
		}
		return &Merkle{Hash: uint32(merkle["hash"].(float64)), Children: children}
	}
	return &Merkle{}
}

func (trie *Merkle) toMapInterface() map[string]interface{} {
	merkleMap := map[string]interface{}{}
	if trie.Hash != 0 {
		merkleMap["hash"] = int32(trie.Hash)
	}

	for key, merkle := range trie.Children {
		if key != "hash" {
			merkleMap[key] = merkle.toMapInterface()
		}
	}

	return merkleMap
}

func (trie *Merkle) ToJSONString() (string, error) {
	merkleMap := trie.toMapInterface()

	jsonString, err := json.Marshal(merkleMap)
	if err != nil {
		return "", err
	}

	return string(jsonString), nil
}

func (trie *Merkle) getKeys() []string {
	j := 0
	keys := make([]string, len(trie.Children))
	for k := range trie.Children {
		keys[j] = k
		j++
	}
	return keys
}

func (trie *Merkle) insertKey(key string, hash uint32) crdt.Merkle {
	if len(key) == 0 {
		return trie
	}

	c := string(key[0])
	cNode := &Merkle{}
	n := trie.Children[c]
	nHash := -1
	if n != nil {
		nHash = int(n.Hash)
		nInserted := n.insertKey(key[1:], hash).(*Merkle)

		cNode.Hash = n.Hash
		for k := range n.Children {
			cNode.Children[k] = n.Children[k]
		}

		if nInserted != nil {
			if nInserted.Hash != 0 {
				cNode.Hash = nInserted.Hash
			}
			for k := range nInserted.Children {
				if nInserted.Children[k] != nil {
					cNode.Children[k] = nInserted.Children[k]
				}
			}
		}

		if hash != 0 {
			cNode.Hash = hash
		}
	}
	if trie.Children == nil {
		trie.Children = map[string]*Merkle{}
	}
	if len(cNode.Children) != 0 {
		trie.Children[c] = cNode
	} else {
		trie.Children[c] = nil
	}
	if nHash != -1 {
		trie.Hash = uint32(nHash) ^ hash
	}
	return trie
}

func (trie *Merkle) Insert(timestamp crdt.Timestamp) {
	hash := timestamp.Hash()
	key := strconv.FormatInt((timestamp.GetMillis() / (1000 * 60)), 3)

	trie.Hash = trie.Hash ^ hash
	trie.insertKey(key, hash)
}

func (trie *Merkle) Prune() crdt.Merkle {
	n := 2

	if trie == nil {
		return nil
	}

	// Checking if empty
	if trie.Children == nil {
		return NewMerkle(trie.Hash)
	}

	keys := trie.getKeys()
	sort.Strings(keys)

	sliceRange := len(keys) - n
	newTrie := NewMerkle(trie.Hash)
	for k := range keys[int(math.Max(0, float64(sliceRange))):] {
		node := trie.Children[keys[k]].Prune()
		if node != nil {
			newTrie.Children[keys[k]] = node.(*Merkle)
		}
	}
	return newTrie
}
