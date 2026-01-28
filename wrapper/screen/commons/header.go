package wrapper_commons

import (
	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/core"
)

type Header struct {
	screen core.Screen
	header []core.Line
}

func NewHeader(screen core.Screen) *Header {
	return &Header{
		screen: screen,
		header: make([]core.Line, 0),
	}
}

func (c *Header) AddHeader(header ...core.Line) *Header {
	c.header = append(c.header, header...)
	return c
}

func (c *Header) ToScreen() core.Screen {
	return core.Screen{
		Name:   c.screen.Name,
		Update: c.Update,
		View:   c.View,
	}
}

func (c *Header) Update(state state.UIState, event core.ScreenEvent) core.ScreenResult {
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
	vm.Header = append(vm.Header, c.header...)
	return vm
}
