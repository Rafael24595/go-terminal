package wrapper_terminal

import (
	"context"
	"fmt"
	"strings"

	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-terminal/engine/model/key"
	"github.com/Rafael24595/go-terminal/engine/model/winsize"
	"github.com/Rafael24595/go-terminal/engine/terminal"
	wrapper_ansi "github.com/Rafael24595/go-terminal/wrapper/ansi"
)

type Console struct {
	context    context.Context
	keyChan    chan key.Key
	resizeChan chan winsize.Winsize
	reader     *inputReader
	buffer     []string
	cursor     uint16
	winsize    winsize.Winsize
	rawmode    uintptr
	color      string
}

func NewConsole() *Console {
	winsize, _ := Size()
	return &Console{
		context: context.Background(),
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

func (c *Console) Context(context context.Context) *Console {
	c.context = context
	return c
}

func (c *Console) ToTerminal() terminal.Terminal {
	return terminal.Terminal{
		OnStart:      c.OnStart,
		OnClose:      c.OnClose,
		ResizeEvents: c.ResizeEvents,
		KeyEvents:    c.KeyEvents,
		Size:         c.Size,
		Clear:        c.Clear,
		Write:        c.Write,
		WriteLine:    c.WriteLine,
		WriteAll:     c.WriteAll,
		Flush:        c.Flush,
	}
}

func (c *Console) OnStart() error {
	rawmode, err := onStart()
	if err != nil {
		return err
	}

	c.rawmode = rawmode
	fmt.Print(c.color + wrapper_ansi.CleanConsole + wrapper_ansi.HideCursor)
	
	return nil
}

func (c *Console) OnClose() error {
	err := onClose(c.rawmode)
	if err != nil {
		return err
	}
	
	fmt.Print(wrapper_ansi.ResetAttrs + wrapper_ansi.CleanConsole + wrapper_ansi.ShowCursor + wrapper_ansi.ResetCursor)
	return nil
}

func (c *Console) ResizeEvents() <-chan winsize.Winsize {
	if c.resizeChan != nil {
		return c.resizeChan
	}

	source := ResizeEvents(c.context)
	c.resizeChan = make(chan winsize.Winsize, 1)
	go c.listenResizeEvents(source)

	return c.resizeChan
}

func (c *Console) listenResizeEvents(source <-chan winsize.Winsize) {
	defer close(c.resizeChan)

	for {
		select {
		case <-c.context.Done():
			return
		case size, ok := <-source:
			if !ok {
				return
			}
			c.defineSize(size)
			select {
			case c.resizeChan <- size:
			default:
			}
		}
	}
}

func (c *Console) KeyEvents() <-chan key.Key {
	if c.keyChan != nil {
		return c.keyChan
	}

	c.keyChan = make(chan key.Key)
	go c.listenKeyEvents()

	return c.keyChan
}

func (c *Console) listenKeyEvents() {
	defer close(c.keyChan)

	for {
		k, err := c.reader.readRune()
		if err != nil {
			return
		}

		select {
		case <-c.context.Done():
			return
		case c.keyChan <- *k:
		}
	}
}

func (c *Console) Size() (winsize.Winsize, error) {
	winsize, err := Size()
	if err != nil {
		return c.winsize, err
	}

	c.buffer = make([]string, winsize.Rows)
	c.cursor = 0
	c.winsize = winsize

	return c.defineSize(winsize), nil
}

func (c *Console) defineSize(winsize winsize.Winsize) winsize.Winsize {
	if c.winsize.Eq(winsize) {
		return c.winsize
	}

	c.buffer = make([]string, winsize.Rows)
	c.cursor = 0
	c.winsize = winsize

	return c.winsize
}

func (c *Console) Clear() error {
	fmt.Print(wrapper_ansi.ResetCursor)
	return nil
}

func (c *Console) Write(fragment string) error {
	newLine := c.buffer[c.cursor] + fragment
	return c.WriteLine(newLine)
}

func (c *Console) WriteLine(line ...string) error {
	for _, l := range line {
		assert.True(winsize.Rows(c.cursor) < c.winsize.Rows,
			"buffer overflow[%d]: the line '%s' cannot be appended at %d position",
			c.winsize.Rows, l, c.cursor,
		)

		c.buffer[c.cursor] = wrapper_ansi.EraseLine + l

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

func (c *Console) clearBuffer() *Console {
	c.cursor = 0
	for i := range c.buffer {
		c.buffer[i] = ""
	}
	return c
}
