package timestamp

// A mutable global clock
var globalClock *Clock

func GetGlobalClock() *Clock {
	return globalClock
}

func SetGlobalClock(clock *Clock) {
	globalClock = clock
}
