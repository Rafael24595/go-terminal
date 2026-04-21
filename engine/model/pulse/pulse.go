package pulse

import (
	"time"

	assert "github.com/Rafael24595/go-assert/assert/runtime"
)

type Pulse struct {
	chn    <-chan time.Time
	tkr    *time.Ticker
	closed bool
}

func New(duration time.Duration) *Pulse {
	return &Pulse{
		chn:    nil,
		tkr:    time.NewTicker(duration),
		closed: false,
	}
}

func (p *Pulse) Listen() <-chan time.Time {
	assert.False(p.closed, "cannot listen a closed pulse")

	return p.chn
}

func (p *Pulse) Enable() *Pulse {
	if p.closed {
		assert.Unreachable("closed pulse cannot be modified")
		return p
	}

	p.chn = p.tkr.C
	return p
}

func (p *Pulse) Disable() *Pulse {
	if p.closed {
		assert.Unreachable("closed pulse cannot be modified")
		return p
	}

	p.chn = nil
	return p
}

func (p *Pulse) Exit() *Pulse {
	if p.closed {
		assert.Unreachable("the pulse is already closed")
		return p
	}

	p.Disable()
	p.tkr.Stop()

	p.closed = true

	return p
}
