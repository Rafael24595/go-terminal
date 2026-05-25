package wrapper_console

import (
	"context"
	"fmt"
	"strings"

	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/terminal"
	"github.com/Rafael24595/go-reacterm-core/wrapper/platform"

	wrapper_ansi "github.com/Rafael24595/go-reacterm-core/wrapper/ansi"
	wrapper_reader "github.com/Rafael24595/go-reacterm-core/wrapper/terminal/reader"
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

func (t *Console) ToTerminal() terminal.Terminal {
	return terminal.Terminal{
		OnStart:      t.OnStart,
		OnClose:      t.OnClose,
		ResizeEvents: t.ResizeEvents,
		KeyEvents:    t.KeyEvents,
		Size:         t.Size,
		Clear:        t.Clear,
		Write:        t.Write,
		WriteAll:     t.WriteAll,
		Flush:        t.Flush,
	}
}

func (t *Console) OnStart() error {
	rawmode, err := platform.OnStart()
	if err != nil {
		return err
	}

	t.rawmode = rawmode

	fmt.Print(t.color + wrapper_ansi.FullReset + wrapper_ansi.HideCursor)

	return nil
}

func (t *Console) OnClose() error {
	err := platform.OnClose(t.rawmode)
	if err != nil {
		return err
	}

	fmt.Print(wrapper_ansi.Reset + wrapper_ansi.FullReset + wrapper_ansi.ShowCursor + wrapper_ansi.CursorHome)

	return nil
}

func (t *Console) ResizeEvents() <-chan winsize.Winsize {
	if t.resizeChan != nil {
		return t.resizeChan
	}

	source := t.strategy(t.context)
	t.resizeChan = make(chan winsize.Winsize, 1)
	go t.listenResizeEvents(source)

	return t.resizeChan
}

func (t *Console) listenResizeEvents(source <-chan winsize.Winsize) {
	defer close(t.resizeChan)

	for {
		select {
		case <-t.context.Done():
			return
		case size, ok := <-source:
			if !ok {
				return
			}
			t.buffer.defineSize(size)
			select {
			case t.resizeChan <- size:
			default:
			}
		}
	}
}

func (t *Console) KeyEvents() <-chan key.Key {
	if t.keyChan != nil {
		return t.keyChan
	}

	t.keyChan = make(chan key.Key)
	go t.listenKeyEvents()

	return t.keyChan
}

func (t *Console) listenKeyEvents() {
	defer close(t.keyChan)

	for {
		k, err := t.reader.ReadKey()
		if err != nil {
			return
		}

		select {
		case <-t.context.Done():
			return
		case t.keyChan <- *k:
		}
	}
}

func (t *Console) Size() (winsize.Winsize, error) {
	winsize, err := platform.Size()
	if err != nil {
		return t.buffer.size, err
	}

	return t.buffer.defineSize(winsize), nil
}

func (t *Console) Clear() error {
	fmt.Print(wrapper_ansi.CursorHome)
	return nil
}

func (t *Console) Write(lines ...string) error {
	for _, l := range lines {
		t.buffer.setLine(wrapper_ansi.ClearLine + l).
			nextLine()
	}
	return nil
}

func (t *Console) WriteAll(text string) error {
	lines := strings.Split(text, "\n")
	return t.Write(lines...)
}

func (t *Console) Flush() error {
	fmt.Print(t.buffer.join("\n"))
	t.buffer.clear()
	return nil
}
