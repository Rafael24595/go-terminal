package wrapper_console

import (
	"strings"

	assert "github.com/Rafael24595/go-assert/assert/runtime"
	"github.com/Rafael24595/go-terminal/engine/model/winsize"
)

type consoleBuffer struct {
	size   winsize.Winsize
	lines  []string
	cursor uint16
}

func newBuffer(size winsize.Winsize) *consoleBuffer {
	return &consoleBuffer{
		size:   size,
		lines:  make([]string, size.Rows),
		cursor: 0,
	}
}

func (c *consoleBuffer) defineSize(size winsize.Winsize) winsize.Winsize {
	if c.size.Eq(size) {
		return c.size
	}

	c.size = size

	c.lines = make([]string, size.Rows)
	c.cursor = 0

	return c.size
}

func (c *consoleBuffer) setLine(line string) *consoleBuffer {
	assert.True(
		winsize.Rows(c.cursor) < c.size.Rows,
		"buffer overflow[%d]: the line '%s' cannot be appended at %d position",
		c.size.Rows, line, c.cursor,
	)

	c.lines[c.cursor] = line
	return c
}

func (c *consoleBuffer) nextLine() bool {
	if c.cursor >= uint16(len(c.lines)) {
		return false
	}

	c.cursor += 1
	return true
}

func (c *consoleBuffer) join(sep string) string {
	return strings.Join(c.lines, sep)
}

func (c *consoleBuffer) clear() *consoleBuffer {
	c.cursor = 0
	for i := range c.lines {
		c.lines[i] = ""
	}
	return c
}
