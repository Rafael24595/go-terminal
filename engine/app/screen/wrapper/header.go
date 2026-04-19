package wrapper

import (
	"github.com/Rafael24595/go-terminal/engine/app/screen"
	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/app/viewmodel"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable/stream/block"
	"github.com/Rafael24595/go-terminal/engine/render/text"
)

type Header struct {
	screen screen.Screen
	header []text.Line
}

func NewHeader(screen screen.Screen) *Header {
	return &Header{
		screen: screen,
		header: make([]text.Line, 0),
	}
}

func (c *Header) AddHeader(header ...text.Line) *Header {
	c.header = append(c.header, header...)
	return c
}

func (c *Header) ToScreen() screen.Screen {
	return screen.Screen{
		Name:       c.screen.Name,
		Definition: c.screen.Definition,
		Update:     c.Update,
		View:       c.View,
		Stack:      c.screen.Stack,
	}
}

func (c *Header) Update(state *state.UIState, event screen.ScreenEvent) screen.ScreenResult {
	result := c.screen.Update(state, event)
	if result.Screen != nil {
		newScreen := NewHeader(*result.Screen).
			AddHeader(c.header...).
			ToScreen()
		result.Screen = &newScreen
	}
	return result
}

func (c *Header) View(state state.UIState) viewmodel.ViewModel {
	vm := c.screen.View(state)

	vm.Header.Unshift(
		block.BlockDrawableFromLines(c.header...),
	)

	return vm
}
