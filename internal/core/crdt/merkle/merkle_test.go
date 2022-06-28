package merkle_test

import (
	"encoding/json"
	"testing"

	"github.com/nathanjisaac/actual-server-go/internal/core/crdt/merkle"
	"github.com/stretchr/testify/assert"
)

func TestMerkle_NewFromMap(t *testing.T) {
	t.Run("parse json string", func(t *testing.T) {
		jsonString := `{
			"1": {
			  "2": {
				"1": {
				  "0": {
					"1": {
					  "0": {
						"0": {
						  "2": {
							"0": {
							  "1": {
								"1": {
								  "0": {
									"2": {
									  "2": {
										"0": {
										  "0": {
											"hash": 1983295247
										  },
										  "hash": 1983295247
										},
										"hash": 1983295247
									  },
									  "hash": 1983295247
									},
									"hash": 1983295247
								  },
								  "hash": 1983295247
								},
								"hash": 1983295247
							  },
							  "hash": 1983295247
							},
							"1": {
							  "0": {
								"1": {
								  "0": {
									"2": {
									  "0": {
										"0": {
										  "0": {
											"hash": 1469038940
										  },
										  "hash": 1469038940
										},
										"hash": 1469038940
									  },
									  "hash": 1469038940
									},
									"hash": 1469038940
								  },
								  "hash": 1469038940
								},
								"hash": 1469038940
							  },
							  "hash": 1469038940
							},
							"hash": 565800531
						  },
						  "hash": 565800531
						},
						"hash": 565800531
					  },
					  "hash": 565800531
					},
					"hash": 565800531
				  },
				  "hash": 565800531
				},
				"hash": 565800531
			  },
			  "hash": 565800531
			},
			"hash": 565800531
		  }`
		var merklemap map[string]interface{}
		err := json.Unmarshal([]byte(jsonString), &merklemap)
		assert.NoError(t, err)

		merklestruct := merkle.NewMerkleFromMap(merklemap)
		assert.Equal(t, uint32(565800531), merklestruct.Hash)
		assert.Equal(t, uint32(565800531), merklestruct.Children["1"].Hash)
		assert.Equal(t, uint32(565800531), merklestruct.Children["1"].Children["2"].Hash)
		assert.Equal(t, uint32(565800531), merklestruct.Children["1"].Children["2"].Children["1"].Hash)
		assert.Equal(t, uint32(565800531), merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].Hash)
		assert.Equal(t, uint32(565800531), merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].Children["1"].Hash)
		assert.Equal(t, uint32(565800531), merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].Children["1"].Children["0"].Hash)
		assert.Equal(t, uint32(565800531), merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].Children["1"].Children["0"].Children["0"].Hash)
		assert.Equal(t, uint32(565800531), merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].Children["1"].Children["0"].Children["0"].Children["2"].Hash)

		assert.Equal(t, uint32(1983295247), merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].Children["1"].Children["0"].Children["0"].Children["2"].Children["0"].Hash)
		assert.Equal(t, uint32(1983295247), merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].Children["1"].Children["0"].Children["0"].Children["2"].Children["0"].Children["1"].Hash)
		assert.Equal(t, uint32(1983295247), merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].Children["1"].Children["0"].Children["0"].Children["2"].Children["0"].Children["1"].Children["1"].Hash)
		assert.Equal(t, uint32(1983295247), merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].Children["1"].Children["0"].Children["0"].Children["2"].Children["0"].Children["1"].Children["1"].Children["0"].Hash)
		assert.Equal(t, uint32(1983295247), merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].Children["1"].Children["0"].Children["0"].Children["2"].Children["0"].Children["1"].Children["1"].Children["0"].Children["2"].Hash)
		assert.Equal(t, uint32(1983295247), merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].Children["1"].Children["0"].Children["0"].Children["2"].Children["0"].Children["1"].Children["1"].Children["0"].Children["2"].Children["2"].Hash)
		assert.Equal(t, uint32(1983295247), merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].Children["1"].Children["0"].Children["0"].Children["2"].Children["0"].Children["1"].Children["1"].Children["0"].Children["2"].Children["2"].Children["0"].Hash)
		assert.Equal(t, uint32(1983295247), merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].Children["1"].Children["0"].Children["0"].Children["2"].Children["0"].Children["1"].Children["1"].Children["0"].Children["2"].Children["2"].Children["0"].Children["0"].Hash)
		assert.Equal(t, 0, len(merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].Children["1"].Children["0"].Children["0"].Children["2"].Children["0"].Children["1"].Children["1"].Children["0"].Children["2"].Children["2"].Children["0"].Children["0"].Children))

		assert.Equal(t, uint32(1469038940), merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].Children["1"].Children["0"].Children["0"].Children["2"].Children["1"].Hash)
		assert.Equal(t, uint32(1469038940), merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].Children["1"].Children["0"].Children["0"].Children["2"].Children["1"].Children["0"].Hash)
		assert.Equal(t, uint32(1469038940), merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].Children["1"].Children["0"].Children["0"].Children["2"].Children["1"].Children["0"].Children["1"].Hash)
		assert.Equal(t, uint32(1469038940), merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].Children["1"].Children["0"].Children["0"].Children["2"].Children["1"].Children["0"].Children["1"].Children["0"].Hash)
		assert.Equal(t, uint32(1469038940), merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].Children["1"].Children["0"].Children["0"].Children["2"].Children["1"].Children["0"].Children["1"].Children["0"].Children["2"].Hash)
		assert.Equal(t, uint32(1469038940), merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].Children["1"].Children["0"].Children["0"].Children["2"].Children["1"].Children["0"].Children["1"].Children["0"].Children["2"].Children["0"].Hash)
		assert.Equal(t, uint32(1469038940), merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].Children["1"].Children["0"].Children["0"].Children["2"].Children["1"].Children["0"].Children["1"].Children["0"].Children["2"].Children["0"].Children["0"].Hash)
		assert.Equal(t, uint32(1469038940), merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].Children["1"].Children["0"].Children["0"].Children["2"].Children["1"].Children["0"].Children["1"].Children["0"].Children["2"].Children["0"].Children["0"].Children["0"].Hash)
		assert.Equal(t, 0, len(merklestruct.Children["1"].Children["2"].Children["1"].Children["0"].Children["1"].Children["0"].Children["0"].Children["2"].Children["1"].Children["0"].Children["1"].Children["0"].Children["2"].Children["0"].Children["0"].Children["0"].Children))
	})
}

func TestMerkle_ToJSONString(t *testing.T) {
	t.Run("parse json string", func(t *testing.T) {
		jsonString := `{"1":{"2":{"1":{"0":{"1":{"0":{"0":{"2":{"0":{"1":{"1":{"0":{"2":{"2":{"0":{"0":{"hash":1983295247},"hash":1983295247},"hash":1983295247},"hash":1983295247},"hash":1983295247},"hash":1983295247},"hash":1983295247},"hash":1983295247},"1":{"0":{"1":{"0":{"2":{"0":{"0":{"0":{"hash":1469038940},"hash":1469038940},"hash":1469038940},"hash":1469038940},"hash":1469038940},"hash":1469038940},"hash":1469038940},"hash":1469038940},"hash":565800531},"hash":565800531},"hash":565800531},"hash":565800531},"hash":565800531},"hash":565800531},"hash":565800531},"hash":565800531},"hash":565800531}`
		var merklemap map[string]interface{}
		err := json.Unmarshal([]byte(jsonString), &merklemap)
		assert.NoError(t, err)

		merklestruct := merkle.NewMerkleFromMap(merklemap)
		output, err := merklestruct.ToJSONString()
		assert.NoError(t, err)
		assert.Equal(t, jsonString, output)
	})
}
