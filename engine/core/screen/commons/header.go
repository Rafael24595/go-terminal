package commons

import (
	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/core/screen"
)

type Header struct {
	screen screen.Screen
	header []core.Line
}

func NewHeader(screen screen.Screen) *Header {
	return &Header{
		screen: screen,
		header: make([]core.Line, 0),
	}
}

func (c *Header) AddHeader(header ...core.Line) *Header {
	c.header = append(c.header, header...)
	return c
}

func (c *Header) ToScreen() screen.Screen {
	return screen.Screen{
		Name:       c.screen.Name,
		Definition: c.screen.Definition,
		Update:     c.Update,
		View:       c.View,
	}
}

func (c *Header) Update(state state.UIState, event screen.ScreenEvent) screen.ScreenResult {
	result := c.screen.Update(state, event)
	if !result.IgnoreParents && result.Screen != nil {
		newScreen := NewHeader(*result.Screen).
			AddHeader(c.header...).
			ToScreen()
		result.Screen = &newScreen
	}
	return result
}

func (c *Header) View(state state.UIState) core.ViewModel {
	vm := c.screen.View(state)
	vm.Header = append(c.header, vm.Header...)
	return vm
}
