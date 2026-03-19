package runtime

import (
	"fmt"
	"math/rand"
	"time"
)

var Instance *Runtime

type Runtime struct {
	sessionId string
	timestamp int64
}

func init() {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	Instance = &Runtime{
		sessionId: newSessionId(rnd),
		timestamp: time.Now().UnixMilli(),
	}
}

func newSessionId(rnd *rand.Rand) string {
	return fmt.Sprintf("%d-%04x", time.Now().UnixNano(), rnd.Uint32())
}

func (r Runtime) SessionId() string {
	return r.sessionId
}

func (r Runtime) Timestamp() int64 {
	return r.timestamp
}
