package wrapper

import (
	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/core/drawable/line"
	"github.com/Rafael24595/go-terminal/engine/core/screen"
	"github.com/Rafael24595/go-terminal/engine/core/text"
)

type Footer struct {
	screen screen.Screen
	footer []text.Line
}

func NewFooter(screen screen.Screen) *Footer {
	return &Footer{
		screen: screen,
		footer: make([]text.Line, 0),
	}
}

func (c *Footer) AddFooter(footer ...text.Line) *Footer {
	c.footer = append(c.footer, footer...)
	return c
}

func (c *Footer) ToScreen() screen.Screen {
	return screen.Screen{
		Name:       c.screen.Name,
		Definition: c.screen.Definition,
		Update:     c.Update,
		View:       c.View,
		Stack:      c.screen.Stack,
	}
}

func (c *Footer) Update(state *state.UIState, event screen.ScreenEvent) screen.ScreenResult {
	result := c.screen.Update(state, event)
	if result.Screen != nil {
		newScreen := NewFooter(*result.Screen).
			AddFooter(c.footer...).
			ToScreen()
		result.Screen = &newScreen
	}
	return result
}

func (c *Footer) View(state state.UIState) core.ViewModel {
	vm := c.screen.View(state)

	vm.Header.Shift(
		line.EagerDrawableFromLines(c.footer...),
	)

	return vm
}
