package wrapper_console

import (
	"strings"

	"github.com/Rafael24595/go-terminal/engine/core/assert"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

const RESET_CURSOR = "\x1b[H"
const CLEAN_CONSOLE = "\x1B[2J\x1B[H"

const HIDE_CURSOR = "\x1b[?25l"
const SHOW_CURSOR = "\x1b[?25h"

type Console struct {
	buffer  []string
	cursor  uint16
	winsize terminal.Winsize
}

func NewConsole() *Console {
	winsize := Size()
	return &Console{
		buffer:  make([]string, winsize.Rows),
		cursor:  0,
		winsize: winsize,
	}
}

func (c *Console) ToTerminal() terminal.Terminal {
	return terminal.Terminal{
		OnStart:   c.OnStart,
		OnClose:   c.OnClose,
		Size:      c.Size,
		Clear:     c.Clear,
		Write:     c.Write,
		WriteLine: c.WriteLine,
		WriteAll:  c.WriteAll,
		Flush:     c.Flush,
	}
}

func (c *Console) OnStart() error {
	print(CLEAN_CONSOLE + HIDE_CURSOR)
	return nil
}

func (c *Console) OnClose() error {
	print(CLEAN_CONSOLE + SHOW_CURSOR + RESET_CURSOR)
	return nil
}

func (c *Console) Size() terminal.Winsize {
	c.winsize = Size()
	return c.winsize
}

func (c *Console) Clear() error {
	print(RESET_CURSOR)
	return nil
}

func (c *Console) Write(fragment string) error {
	newLine := c.buffer[c.cursor] + fragment
	return c.WriteLine(newLine)
}

func (c *Console) WriteLine(line ...string) error {
	for _, l := range line {
		assert.AssertfTrue(
			c.cursor < c.winsize.Rows,
			"buffer overflow[%d]: the line '%s' cannot be appended at %d position",
			c.winsize.Rows, l, c.cursor,
		)

		// assert.AssertfFalse(len(l) > int(c.winsize.Cols),
		// 	"line overflow[%d]: the line '%s' has lenght %d",
		// 	c.winsize.Cols, l, len(l),
		// )

		c.buffer[c.cursor] = l

		c.cursor += 1
	}
	return nil
}

func (c *Console) WriteAll(text string) error {
	return c.WriteLine(strings.Split(text, "\n")...)
}

func (c *Console) Flush() error {
	print(strings.Join(c.buffer, "\n"))
	c.clearBuffer()
	return nil
}

func (c *Console) clearBuffer() error {
	c.cursor = 0
	c.buffer = make([]string, c.winsize.Rows)
	return nil
}
