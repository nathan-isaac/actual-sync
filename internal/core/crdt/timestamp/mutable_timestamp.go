package timestamp

type MutableTimestamp struct {
	Timestamp
}

func NewMutableTimestamp(millis int64, counter int64, node string) *MutableTimestamp {
	return &MutableTimestamp{Timestamp{millis: millis, counter: counter, node: node}}
}

func NewMutableTimestampFrom(ts Timestamp) *MutableTimestamp {
	return &MutableTimestamp{ts}
}

func (mts *MutableTimestamp) SetMillis(millis int64) {
	mts.millis = millis
}

func (mts *MutableTimestamp) SetCounter(counter int64) {
	mts.counter = counter
}

func (mts *MutableTimestamp) SetNode(node string) {
	mts.node = node
}
