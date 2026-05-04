package clock

import (
	"sync/atomic"
	"time"
)

var globalCounter int64

type Clock func() int64

func UnixMilliClock() int64 {
	return time.Now().UnixMilli()
}

func GlobalCounterClock() int64 {
	return atomic.AddInt64(&globalCounter, 1)

}
