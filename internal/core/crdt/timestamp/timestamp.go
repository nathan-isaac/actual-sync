package timestamp

import (
	"strconv"
	"strings"
	"time"

	internal_errors "github.com/nathanjisaac/actual-server-go/internal/errors"
	"github.com/spaolacci/murmur3"
)

type Timestamp struct {
	millis  int64
	counter int64
	node    string
}

const timestampFormatISO string = "2006-01-02T15:04:05.000Z07:00"

func NewTimestamp(millis int64, counter int64, node string) *Timestamp {
	return &Timestamp{millis: millis, counter: counter, node: node}
}

func ParseTimestamp(str string) (*Timestamp, error) {
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
		return &Timestamp{millis: millis, counter: counter, node: node}, nil
	}
	return nil, internal_errors.ErrTimestampUnableToParse
}

func (ts *Timestamp) ToString() string {
	isoString := time.UnixMilli(ts.millis).UTC().Format(timestampFormatISO)
	cntrString := "0000" + strings.ToUpper(strconv.FormatInt(ts.counter, 16))
	nodeString := "0000000000000000" + ts.node
	str := []string{isoString, cntrString[len(cntrString)-4:], nodeString[len(nodeString)-16:]}
	return strings.Join(str, "-")
}

func (ts *Timestamp) GetMillis() int64 {
	return ts.millis
}

func (ts *Timestamp) GetCounter() int64 {
	return ts.counter
}

func (ts *Timestamp) GetNode() string {
	return ts.node
}

func (ts *Timestamp) Hash() uint32 {
	return murmur3.Sum32([]byte(ts.ToString()))
}
