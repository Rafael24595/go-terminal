package primitive

import (
	"testing"

	"github.com/Rafael24595/go-terminal/engine/core/style"
	"github.com/Rafael24595/go-terminal/test/support/assert"
	"github.com/Rafael24595/go-terminal/test/support/mock"
)

func TestCursor_SelectionLogic(t *testing.T) {
	c := NewCursor(true)
	buff := []rune("Golang")

	c.MoveSelectTo(buff, 3, 1)
	assert.Equal(t, c.SelectStart(), 1)
	assert.Equal(t, c.SelectEnd(), 3)

	c.MoveSelectTo(buff, 1, 3)
	assert.Equal(t, c.SelectStart(), 1)
	assert.Equal(t, c.SelectEnd(), 3)

	c.MoveCaretTo(buff, 99)
	assert.Equal(t, c.Caret(), 6)
}

func TestCursor_BlinkingLogic(t *testing.T) {
	clock := &mock.TestClock{Time: 0}

	c := NewCursor(true)
	c.clock = clock.Now

	clock.Advance(blink_ms + 1)

	assert.Equal(t, c.BlinkStyle(), style.AtmSelect)

	clock.Advance(blink_ms + 1)

	assert.Equal(t, c.BlinkStyle(), style.AtmNone)

	clock.Advance(blink_ms + 1)
	assert.Equal(t, c.BlinkStyle(), style.AtmSelect)
}
