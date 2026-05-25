package wrapper_console

import (
	"strings"

	assert "github.com/Rafael24595/go-assert/assert/runtime"
	
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
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

func (b *consoleBuffer) defineSize(size winsize.Winsize) winsize.Winsize {
	if b.size.Eq(size) {
		return b.size
	}

	b.size = size

	b.lines = make([]string, size.Rows)
	b.cursor = 0

	return b.size
}

func (b *consoleBuffer) setLine(line string) *consoleBuffer {
	assert.True(
		winsize.Rows(b.cursor) < b.size.Rows,
		"buffer overflow[%d]: the line '%s' cannot be appended at %d position",
		b.size.Rows, line, b.cursor,
	)

	b.lines[b.cursor] = line
	return b
}

func (b *consoleBuffer) nextLine() bool {
	if b.cursor >= uint16(len(b.lines)) {
		return false
	}

	b.cursor += 1
	return true
}

func (b *consoleBuffer) join(sep string) string {
	return strings.Join(b.lines, sep)
}

func (b *consoleBuffer) clear() *consoleBuffer {
	b.cursor = 0
	for i := range b.lines {
		b.lines[i] = ""
	}
	return b
}
