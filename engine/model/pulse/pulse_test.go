package pulse

import (
	"testing"
	"time"

	assert "github.com/Rafael24595/go-assert/assert/test"
)

func TestPulse_Toggle(t *testing.T) {
	p := New(50 * time.Millisecond)
	defer p.Exit()

	assert.Nil(t, p.Listen())

	p.Enable()
	assert.NotNil(t, p.Listen())

	p.Disable()
	assert.Nil(t, p.Listen())
}

func TestPulse_Reception(t *testing.T) {
    p := New(10 * time.Millisecond)
    defer p.Exit()
    
    p.Enable()

    select {
    case <-p.Listen():
    case <-time.After(100 * time.Millisecond):
        t.Error("timeout: the pulse never came")
    }
}

func TestPulse_Exit(t *testing.T) {
    p := New(10 * time.Millisecond)
	
    p.Enable()
    p.Exit()

	assert.Panic(t, func() {
		p.Enable()
	})

	assert.Panic(t, func() {
		p.Disable()
	})

	assert.Panic(t, func() {
		p.Listen()
	})

	assert.Panic(t, func() {
		p.Exit()
	})
}
