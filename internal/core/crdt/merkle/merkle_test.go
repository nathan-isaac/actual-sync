package merkle_test

import (
	"encoding/json"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/nathanjisaac/actual-server-go/internal/core/crdt"
	"github.com/nathanjisaac/actual-server-go/internal/core/crdt/merkle"
	"github.com/nathanjisaac/actual-server-go/internal/core/crdt/timestamp"
	internal_errors "github.com/nathanjisaac/actual-server-go/internal/errors"
	"github.com/stretchr/testify/assert"
)

func TestMerkle_NewFromMap(t *testing.T) {
	t.Run("parse json string", func(t *testing.T) {
		jsonString := `{"1":{"2":{"1":{"0":{"1":{"0":{"0":{"2":{"0":{"1":{"1":{"0":{"2":{"2":{"0":{"0":` +
			`{"hash":1983295247},"hash":1983295247},"hash":1983295247},"hash":1983295247},"hash":1983295247},` +
			`"hash":1983295247},"hash":1983295247},"hash":1983295247},"1":{"0":{"1":{"0":{"2":{"0":{"0":{"0":` +
			`{"hash":1469038940},"hash":1469038940},"hash":1469038940},"hash":1469038940},"hash":1469038940},` +
			`"hash":1469038940},"hash":1469038940},"hash":1469038940},"hash":565800531},"hash":565800531},` +
			`"hash":565800531},"hash":565800531},"hash":565800531},"hash":565800531},"hash":565800531},` +
			`"hash":565800531},"hash":565800531}`
		var merklemap map[string]interface{}
		err := json.Unmarshal([]byte(jsonString), &merklemap)
		assert.NoError(t, err)

		merklestruct := merkle.NewMerkleFromMap(merklemap)
		assert.Equal(t, uint32(565800531), merklestruct.Hash)
		assert.Equal(t, uint32(565800531), merklestruct.Children["1"].Hash)
		assert.Equal(t, uint32(565800531), merklestruct.Children["1"].Children["2"].Hash)
		assert.Equal(t, uint32(565800531), merklestruct.Children["1"].Children["2"].Children["1"].Hash)
		assert.Equal(t, uint32(565800531), merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].Hash)
		assert.Equal(t, uint32(565800531), merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].
			Children["1"].Hash)
		assert.Equal(t, uint32(565800531), merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].
			Children["1"].Children["0"].Hash)
		assert.Equal(t, uint32(565800531), merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].
			Children["1"].Children["0"].Children["0"].Hash)
		assert.Equal(t, uint32(565800531), merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].
			Children["1"].Children["0"].Children["0"].Children["2"].Hash)

		assert.Equal(t, uint32(1983295247), merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].
			Children["1"].Children["0"].Children["0"].Children["2"].Children["0"].Hash)
		assert.Equal(t, uint32(1983295247), merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].
			Children["1"].Children["0"].Children["0"].Children["2"].Children["0"].Children["1"].Hash)
		assert.Equal(t, uint32(1983295247), merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].
			Children["1"].Children["0"].Children["0"].Children["2"].Children["0"].Children["1"].Children["1"].Hash)
		assert.Equal(t, uint32(1983295247), merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].
			Children["1"].Children["0"].Children["0"].Children["2"].Children["0"].Children["1"].Children["1"].Children["0"].Hash)
		assert.Equal(t, uint32(1983295247), merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].
			Children["1"].Children["0"].Children["0"].Children["2"].Children["0"].Children["1"].Children["1"].Children["0"].
			Children["2"].Hash)
		assert.Equal(t, uint32(1983295247), merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].
			Children["1"].Children["0"].Children["0"].Children["2"].Children["0"].Children["1"].Children["1"].Children["0"].
			Children["2"].Children["2"].Hash)
		assert.Equal(t, uint32(1983295247), merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].
			Children["1"].Children["0"].Children["0"].Children["2"].Children["0"].Children["1"].Children["1"].Children["0"].
			Children["2"].Children["2"].Children["0"].Hash)
		assert.Equal(t, uint32(1983295247), merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].
			Children["1"].Children["0"].Children["0"].Children["2"].Children["0"].Children["1"].Children["1"].Children["0"].
			Children["2"].Children["2"].Children["0"].Children["0"].Hash)
		assert.Equal(t, 0, len(merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].Children["1"].
			Children["0"].Children["0"].Children["2"].Children["0"].Children["1"].Children["1"].Children["0"].Children["2"].
			Children["2"].Children["0"].Children["0"].Children))

		assert.Equal(t, uint32(1469038940), merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].
			Children["1"].Children["0"].Children["0"].Children["2"].Children["1"].Hash)
		assert.Equal(t, uint32(1469038940), merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].
			Children["1"].Children["0"].Children["0"].Children["2"].Children["1"].Children["0"].Hash)
		assert.Equal(t, uint32(1469038940), merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].
			Children["1"].Children["0"].Children["0"].Children["2"].Children["1"].Children["0"].Children["1"].Hash)
		assert.Equal(t, uint32(1469038940), merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].
			Children["1"].Children["0"].Children["0"].Children["2"].Children["1"].Children["0"].Children["1"].Children["0"].Hash)
		assert.Equal(t, uint32(1469038940), merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].
			Children["1"].Children["0"].Children["0"].Children["2"].Children["1"].Children["0"].Children["1"].Children["0"].
			Children["2"].Hash)
		assert.Equal(t, uint32(1469038940), merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].
			Children["1"].Children["0"].Children["0"].Children["2"].Children["1"].Children["0"].Children["1"].Children["0"].
			Children["2"].Children["0"].Hash)
		assert.Equal(t, uint32(1469038940), merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].
			Children["1"].Children["0"].Children["0"].Children["2"].Children["1"].Children["0"].Children["1"].Children["0"].
			Children["2"].Children["0"].Children["0"].Hash)
		assert.Equal(t, uint32(1469038940), merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].
			Children["1"].Children["0"].Children["0"].Children["2"].Children["1"].Children["0"].Children["1"].Children["0"].
			Children["2"].Children["0"].Children["0"].Children["0"].Hash)
		assert.Equal(t, 0, len(merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].Children["1"].
			Children["0"].Children["0"].Children["2"].Children["1"].Children["0"].Children["1"].Children["0"].Children["2"].
			Children["0"].Children["0"].Children["0"].Children))
	})
}

func TestMerkle_ToJSONString(t *testing.T) {
	t.Run("return json strung from merkle struct", func(t *testing.T) {
		jsonString := `{"1":{"2":{"1":{"0":{"1":{"0":{"0":{"2":{"0":{"1":{"1":{"0":{"2":{"2":{"0":{"0":` +
			`{"hash":1983295247},"hash":1983295247},"hash":1983295247},"hash":1983295247},"hash":1983295247},` +
			`"hash":1983295247},"hash":1983295247},"hash":1983295247},"1":{"0":{"1":{"0":{"2":{"0":{"0":{"0":` +
			`{"hash":1469038940},"hash":1469038940},"hash":1469038940},"hash":1469038940},"hash":1469038940},` +
			`"hash":1469038940},"hash":1469038940},"hash":1469038940},"hash":565800531},"hash":565800531},` +
			`"hash":565800531},"hash":565800531},"hash":565800531},"hash":565800531},"hash":565800531},` +
			`"hash":565800531},"hash":565800531}`
		var merklemap map[string]interface{}
		err := json.Unmarshal([]byte(jsonString), &merklemap)
		assert.NoError(t, err)

		merklestruct := merkle.NewMerkleFromMap(merklemap)
		output, err := merklestruct.ToJSONString()
		assert.NoError(t, err)
		assert.Equal(t, jsonString, output)
	})
}

func TestMerkle_Insert(t *testing.T) {
	t.Run("adding an item works", func(t *testing.T) {
		merklestruct := merkle.NewMerkle(0)

		ts, err := timestamp.ParseTimestamp("2018-11-12T13:21:40.122Z-0000-0123456789ABCDEF")
		assert.NoError(t, err)
		merklestruct.Insert(ts)

		ts, err = timestamp.ParseTimestamp("2018-11-13T13:21:40.122Z-0000-0123456789ABCDEF")
		assert.NoError(t, err)
		merklestruct.Insert(ts)

		assert.Equal(t, uint32(565800531), merklestruct.Hash)
		assert.Equal(t, uint32(565800531), merklestruct.Children["1"].Hash)
		assert.Equal(t, uint32(565800531), merklestruct.Children["1"].Children["2"].Hash)
		assert.Equal(t, uint32(565800531), merklestruct.Children["1"].Children["2"].Children["1"].Hash)
		assert.Equal(t, uint32(565800531), merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].Hash)
		assert.Equal(t, uint32(565800531), merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].
			Children["1"].Hash)
		assert.Equal(t, uint32(565800531), merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].
			Children["1"].Children["0"].Hash)
		assert.Equal(t, uint32(565800531), merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].
			Children["1"].Children["0"].Children["0"].Hash)
		assert.Equal(t, uint32(565800531), merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].
			Children["1"].Children["0"].Children["0"].Children["2"].Hash)

		assert.Equal(t, uint32(1983295247), merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].
			Children["1"].Children["0"].Children["0"].Children["2"].Children["0"].Hash)
		assert.Equal(t, uint32(1983295247), merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].
			Children["1"].Children["0"].Children["0"].Children["2"].Children["0"].Children["1"].Hash)
		assert.Equal(t, uint32(1983295247), merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].
			Children["1"].Children["0"].Children["0"].Children["2"].Children["0"].Children["1"].Children["1"].Hash)
		assert.Equal(t, uint32(1983295247), merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].
			Children["1"].Children["0"].Children["0"].Children["2"].Children["0"].Children["1"].Children["1"].Children["0"].Hash)
		assert.Equal(t, uint32(1983295247), merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].
			Children["1"].Children["0"].Children["0"].Children["2"].Children["0"].Children["1"].Children["1"].Children["0"].
			Children["2"].Hash)
		assert.Equal(t, uint32(1983295247), merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].
			Children["1"].Children["0"].Children["0"].Children["2"].Children["0"].Children["1"].Children["1"].Children["0"].
			Children["2"].Children["2"].Hash)
		assert.Equal(t, uint32(1983295247), merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].
			Children["1"].Children["0"].Children["0"].Children["2"].Children["0"].Children["1"].Children["1"].Children["0"].
			Children["2"].Children["2"].Children["0"].Hash)
		assert.Equal(t, uint32(1983295247), merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].
			Children["1"].Children["0"].Children["0"].Children["2"].Children["0"].Children["1"].Children["1"].Children["0"].
			Children["2"].Children["2"].Children["0"].Children["0"].Hash)
		assert.Equal(t, 0, len(merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].Children["1"].
			Children["0"].Children["0"].Children["2"].Children["0"].Children["1"].Children["1"].Children["0"].Children["2"].
			Children["2"].Children["0"].Children["0"].Children))

		assert.Equal(t, uint32(1469038940), merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].
			Children["1"].Children["0"].Children["0"].Children["2"].Children["1"].Hash)
		assert.Equal(t, uint32(1469038940), merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].
			Children["1"].Children["0"].Children["0"].Children["2"].Children["1"].Children["0"].Hash)
		assert.Equal(t, uint32(1469038940), merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].
			Children["1"].Children["0"].Children["0"].Children["2"].Children["1"].Children["0"].Children["1"].Hash)
		assert.Equal(t, uint32(1469038940), merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].
			Children["1"].Children["0"].Children["0"].Children["2"].Children["1"].Children["0"].Children["1"].Children["0"].Hash)
		assert.Equal(t, uint32(1469038940), merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].
			Children["1"].Children["0"].Children["0"].Children["2"].Children["1"].Children["0"].Children["1"].Children["0"].
			Children["2"].Hash)
		assert.Equal(t, uint32(1469038940), merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].
			Children["1"].Children["0"].Children["0"].Children["2"].Children["1"].Children["0"].Children["1"].Children["0"].
			Children["2"].Children["0"].Hash)
		assert.Equal(t, uint32(1469038940), merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].
			Children["1"].Children["0"].Children["0"].Children["2"].Children["1"].Children["0"].Children["1"].Children["0"].
			Children["2"].Children["0"].Children["0"].Hash)
		assert.Equal(t, uint32(1469038940), merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].
			Children["1"].Children["0"].Children["0"].Children["2"].Children["1"].Children["0"].Children["1"].Children["0"].
			Children["2"].Children["0"].Children["0"].Children["0"].Hash)
		assert.Equal(t, 0, len(merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].Children["1"].
			Children["0"].Children["0"].Children["2"].Children["1"].Children["0"].Children["1"].Children["0"].Children["2"].
			Children["0"].Children["0"].Children["0"].Children))
	})
}

type testStubTimestamp struct {
	millis  int64
	counter int64
	node    string
	hash    uint32
}

func (ts *testStubTimestamp) ToString() string {
	return ""
}

func (ts *testStubTimestamp) GetMillis() int64 {
	return ts.millis
}

func (ts *testStubTimestamp) GetCounter() int64 {
	return ts.counter
}

func (ts *testStubTimestamp) GetNode() string {
	return ts.node
}

func (ts *testStubTimestamp) Hash() uint32 {
	return ts.hash
}

func parseTimestampStub(str string, hash uint32) (*testStubTimestamp, error) {
	parts := strings.Split(str, "-")
	if parts != nil && (len(parts) == 5) {
		t, err := time.Parse(time.RFC3339, strings.Join(parts[:3], "-"))
		if err != nil {
			return nil, err
		}
		millis := t.UnixMilli()
		counter, err := strconv.ParseInt(parts[3], 16, 64)
		if err != nil {
			return nil, err
		}
		node := parts[4]
		return &testStubTimestamp{millis: millis, counter: counter, node: node, hash: hash}, nil
	}
	return nil, internal_errors.ErrTimestampUnableToParse
}

func TestMerkle_Prune(t *testing.T) {
	t.Run("pruning works and keeps correct hashes", func(t *testing.T) {
		times := make([]crdt.Timestamp, 12)
		ts, err := parseTimestampStub("2018-11-01T01:00:00.000Z-0000-0123456789ABCDEF", uint32(1000))
		assert.NoError(t, err)
		times[0] = ts
		ts, err = parseTimestampStub("2018-11-01T01:09:00.000Z-0000-0123456789ABCDEF", uint32(1100))
		assert.NoError(t, err)
		times[1] = ts
		ts, err = parseTimestampStub("2018-11-01T01:18:00.000Z-0000-0123456789ABCDEF", uint32(1200))
		assert.NoError(t, err)
		times[2] = ts
		ts, err = parseTimestampStub("2018-11-01T01:27:00.000Z-0000-0123456789ABCDEF", uint32(1300))
		assert.NoError(t, err)
		times[3] = ts
		ts, err = parseTimestampStub("2018-11-01T01:36:00.000Z-0000-0123456789ABCDEF", uint32(1400))
		assert.NoError(t, err)
		times[4] = ts
		ts, err = parseTimestampStub("2018-11-01T01:45:00.000Z-0000-0123456789ABCDEF", uint32(1500))
		assert.NoError(t, err)
		times[5] = ts
		ts, err = parseTimestampStub("2018-11-01T01:54:00.000Z-0000-0123456789ABCDEF", uint32(1600))
		assert.NoError(t, err)
		times[6] = ts
		ts, err = parseTimestampStub("2018-11-01T02:03:00.000Z-0000-0123456789ABCDEF", uint32(1700))
		assert.NoError(t, err)
		times[7] = ts
		ts, err = parseTimestampStub("2018-11-01T02:10:00.000Z-0000-0123456789ABCDEF", uint32(1800))
		assert.NoError(t, err)
		times[8] = ts
		ts, err = parseTimestampStub("2018-11-01T02:19:00.000Z-0000-0123456789ABCDEF", uint32(1900))
		assert.NoError(t, err)
		times[9] = ts
		ts, err = parseTimestampStub("2018-11-01T02:28:00.000Z-0000-0123456789ABCDEF", uint32(2000))
		assert.NoError(t, err)
		times[10] = ts
		ts, err = parseTimestampStub("2018-11-01T02:37:00.000Z-0000-0123456789ABCDEF", uint32(2100))
		assert.NoError(t, err)
		times[11] = ts

		merklestruct := merkle.NewMerkle(0)
		for _, ts := range times {
			merklestruct.Insert(ts)
		}
		jsonString := `{"1":{"2":{"1":{"0":{"0":{"2":{"2":{"2":{"1":{"2":{"2":{"0":{"0":{"1":{"2":{"0":{"hash":1000},` +
			`"hash":1000},"hash":1000},"2":{"2":{"0":{"hash":1100},"hash":1100},"hash":1100},"hash":1956},"1":{"0":` +
			`{"2":{"0":{"hash":1200},"hash":1200},"hash":1200},"1":{"2":{"0":{"hash":1300},"hash":1300},"hash":1300},` +
			`"2":{"2":{"0":{"hash":1400},"hash":1400},"hash":1400},"hash":1244},"2":{"0":{"2":{"0":{"hash":1500},"hash":1500},` +
			`"hash":1500},"1":{"2":{"0":{"hash":1600},"hash":1600},"hash":1600},"2":{"2":{"0":{"hash":1700},"hash":1700},` +
			`"hash":1700},"hash":1336},"hash":1600},"1":{"0":{"0":{"1":{"1":{"hash":1800},"hash":1800},"hash":1800},` +
			`"1":{"1":{"1":{"hash":1900},"hash":1900},"hash":1900},"2":{"1":{"1":{"hash":2000},"hash":2000},"hash":2000},` +
			`"hash":1972},"1":{"0":{"1":{"1":{"hash":2100},"hash":2100},"hash":2100},"hash":2100},"hash":3968},"hash":2496},` +
			`"hash":2496},"hash":2496},"hash":2496},"hash":2496},"hash":2496},"hash":2496},"hash":2496},"hash":2496},` +
			`"hash":2496},"hash":2496},"hash":2496}`
		jsonOutput, err := merklestruct.ToJSONString()
		assert.NoError(t, err)

		assert.Equal(t, uint32(2496), merklestruct.Hash)
		assert.Equal(t, jsonString, jsonOutput)

		pruned, ok := merklestruct.Prune().(*merkle.Merkle)
		assert.Equal(t, true, ok)
		jsonString = `{"1":{"2":{"1":{"0":{"0":{"2":{"2":{"2":{"1":{"2":{"2":{"0":{"1":{"1":{"2":{"0":{"hash":1300},` +
			`"hash":1300},"hash":1300},"2":{"2":{"0":{"hash":1400},"hash":1400},"hash":1400},"hash":1244},"2":{"1":` +
			`{"2":{"0":{"hash":1600},"hash":1600},"hash":1600},"2":{"2":{"0":{"hash":1700},"hash":1700},"hash":1700},` +
			`"hash":1336},"hash":1600},"1":{"0":{"1":{"1":{"1":{"hash":1900},"hash":1900},"hash":1900},"2":{"1":{"1":` +
			`{"hash":2000},"hash":2000},"hash":2000},"hash":1972},"1":{"0":{"1":{"1":{"hash":2100},"hash":2100},"hash":2100},` +
			`"hash":2100},"hash":3968},"hash":2496},"hash":2496},"hash":2496},"hash":2496},"hash":2496},"hash":2496},` +
			`"hash":2496},"hash":2496},"hash":2496},"hash":2496},"hash":2496},"hash":2496}`
		jsonOutput, err = pruned.ToJSONString()
		assert.NoError(t, err)

		assert.Equal(t, uint32(2496), pruned.Hash)
		assert.Equal(t, jsonString, jsonOutput)
	})
}
