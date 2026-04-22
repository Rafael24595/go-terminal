package wrapper_console

import (
	"context"
	"fmt"
	"strings"

	"github.com/Rafael24595/go-terminal/engine/model/key"
	"github.com/Rafael24595/go-terminal/engine/model/winsize"
	"github.com/Rafael24595/go-terminal/engine/terminal"
	"github.com/Rafael24595/go-terminal/wrapper/platform"

	wrapper_ansi "github.com/Rafael24595/go-terminal/wrapper/ansi"
	wrapper_reader "github.com/Rafael24595/go-terminal/wrapper/terminal/reader"
)

type Console struct {
	context    context.Context
	strategy   resizeStrategy
	keyChan    chan key.Key
	resizeChan chan winsize.Winsize
	reader     *wrapper_reader.KeyReader
	buffer     *consoleBuffer
	rawmode    uintptr
	color      string
}

func newConsole() *Console {
	winsize, _ := platform.Size()
	return &Console{
		context:  context.Background(),
		strategy: defaultStrategy(),
		reader:   wrapper_reader.New(),
		buffer:   newBuffer(winsize),
		color:    "",
	}
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
		WriteAll:     c.WriteAll,
		Flush:        c.Flush,
	}
}

func (c *Console) OnStart() error {
	rawmode, err := platform.OnStart()
	if err != nil {
		return err
	}

	c.rawmode = rawmode

	fmt.Print(c.color + wrapper_ansi.FullReset + wrapper_ansi.HideCursor)

	return nil
}

func (c *Console) OnClose() error {
	err := platform.OnClose(c.rawmode)
	if err != nil {
		return err
	}

	fmt.Print(wrapper_ansi.Reset + wrapper_ansi.FullReset + wrapper_ansi.ShowCursor + wrapper_ansi.CursorHome)

	return nil
}

func (c *Console) ResizeEvents() <-chan winsize.Winsize {
	if c.resizeChan != nil {
		return c.resizeChan
	}

	source := c.strategy(c.context)
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
			c.buffer.defineSize(size)
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
		k, err := c.reader.ReadKey()
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
	winsize, err := platform.Size()
	if err != nil {
		return c.buffer.size, err
	}

	return c.buffer.defineSize(winsize), nil
}

func (c *Console) Clear() error {
	fmt.Print(wrapper_ansi.CursorHome)
	return nil
}

func (c *Console) Write(lines ...string) error {
	for _, l := range lines {
		c.buffer.setLine(wrapper_ansi.ClearLine + l).
			nextLine()
	}
	return nil
}

func (c *Console) WriteAll(text string) error {
	lines := strings.Split(text, "\n")
	return c.Write(lines...)
}

func (c *Console) Flush() error {
	fmt.Print(c.buffer.join("\n"))
	c.buffer.clear()
	return nil
}
