package input

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	
	"github.com/Rafael24595/go-terminal/engine/render/style"
	"github.com/Rafael24595/go-terminal/test/support/mock"
)

func TestCursor_SelectionLogic(t *testing.T) {
	c := NewTextCursor(true)
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

	c := NewTextCursor(true)
	c.clock = clock.Now

	clock.Advance(blink_ms + 1)

	assert.Equal(t, c.BlinkStyle(), style.AtmSelect)

	clock.Advance(blink_ms + 1)

	assert.Equal(t, c.BlinkStyle(), style.AtmNone)

	clock.Advance(blink_ms + 1)
	assert.Equal(t, c.BlinkStyle(), style.AtmSelect)
}
