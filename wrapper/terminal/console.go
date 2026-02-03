package wrapper_terminal

import (
	"fmt"
	"strings"

	"github.com/Rafael24595/go-terminal/engine/core/assert"
	"github.com/Rafael24595/go-terminal/engine/core/key"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

const RESET_ATTRS = "\x1b[0m"
const RESET_CURSOR = "\x1b[H"
const CLEAN_CONSOLE = "\x1B[2J\x1B[H"

const ERASE_LINE = "\r\033[K"

const HIDE_CURSOR = "\x1b[?25l"
const SHOW_CURSOR = "\x1b[?25h"

type Console struct {
	reader  *inputReader
	buffer  []string
	cursor  uint16
	winsize terminal.Winsize
	rawmode uintptr
	color   string
}

func NewConsole() *Console {
	winsize := Size()
	return &Console{
		reader:  newInputReader(),
		buffer:  make([]string, winsize.Rows),
		cursor:  0,
		winsize: winsize,
		color:   "",
	}
}

func (c *Console) Color(color string) *Console {
	c.color = color
	return c
}

func (c *Console) Update() *Console {
	winsize := Size()
	if winsize.Cols == c.winsize.Cols && winsize.Rows == c.winsize.Rows {
		return c
	}

	c.buffer = make([]string, winsize.Rows)
	c.cursor = 0
	c.winsize = winsize

	return c
}

func (c *Console) ToTerminal() terminal.Terminal {
	return terminal.Terminal{
		OnStart:   c.OnStart,
		OnClose:   c.OnClose,
		Size:      c.Size,
		Clear:     c.Clear,
		ReadKey:   c.ReadKey,
		Write:     c.Write,
		WriteLine: c.WriteLine,
		WriteAll:  c.WriteAll,
		Flush:     c.Flush,
	}
}

func (c *Console) OnStart() error {
	c.rawmode, _ = onStart()
	fmt.Print(c.color + CLEAN_CONSOLE + HIDE_CURSOR)
	return nil
}

func (c *Console) OnClose() error {
	onClose(c.rawmode)
	fmt.Print(RESET_ATTRS + CLEAN_CONSOLE + SHOW_CURSOR + RESET_CURSOR)
	return nil
}

func (c *Console) Size() terminal.Winsize {
	return Size()
}

func (c *Console) Clear() error {
	fmt.Print(RESET_CURSOR)
	return nil
}

func (c *Console) ReadKey() (*key.Key, error) {
	return c.reader.readRune()
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

		// assert.AssertfFalse(utf8.RuneCountInString(l) > int(c.winsize.Cols),
		// 	"line overflow[%d]: the line '%s' has lenght %d",
		// 	c.winsize.Cols, l, utf8.RuneCountInString(l),
		// )

		c.buffer[c.cursor] = ERASE_LINE + l

		c.cursor += 1
	}
	return nil
}

func (c *Console) WriteAll(text string) error {
	return c.WriteLine(strings.Split(text, "\n")...)
}

func (c *Console) Flush() error {
	fmt.Print(strings.Join(c.buffer, "\n"))
	c.clearBuffer()
	return nil
}

func (c *Console) clearBuffer() error {
	c.cursor = 0
	c.buffer = make([]string, c.winsize.Rows)
	return nil
}
