package timestamp_test

import (
	"testing"

	"github.com/nathanjisaac/actual-server-go/internal/core/crdt/timestamp"
	"github.com/stretchr/testify/assert"
)

func TestTimestamp_Parse(t *testing.T) {
	t.Run("should not parse", func(t *testing.T) {
		invalidInputs := []string{
			"",
			" ",
			"0",
			"invalid",
			"1969-1-1T0:0:0.0Z-0-0-0",
		}

		for i := range invalidInputs {
			result, err := timestamp.ParseTimestamp(invalidInputs[i])
			assert.Nil(t, result)
			assert.Error(t, err)
		}
	})

	t.Run("should parse", func(t *testing.T) {
		validInputs := []string{
			"1970-01-01T00:00:00.000Z-0000-0000000000000000",
			"2015-04-24T22:23:42.123Z-1000-0123456789ABCDEF",
			"9999-12-31T23:59:59.999Z-FFFF-FFFFFFFFFFFFFFFF",
		}

		for i := range validInputs {
			result, err := timestamp.ParseTimestamp(validInputs[i])
			assert.NoError(t, err)
			assert.GreaterOrEqual(t, result.GetMillis(), int64(0))
			assert.Less(t, result.GetMillis(), int64(253402300800000))
			assert.GreaterOrEqual(t, result.GetCounter(), int64(0))
			assert.Less(t, result.GetCounter(), int64(65536))
			assert.Equal(t, validInputs[i], result.ToString())
		}
	})
}
