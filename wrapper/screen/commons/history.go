package wrapper_commons

import (
	"fmt"

	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/core"
)

type History struct {
	history *core.Screen
	screen core.Screen
}

func NewHistory(screen core.Screen) *History {
	return &History{
		screen: screen,
	}
}

func (c *History) ToScreen() core.Screen {
	return core.Screen{
		Name:   c.screen.Name,
		Update: c.update,
		View:   c.view,
	}
}

func (c *History) update(state state.UIState, event core.ScreenEvent) core.ScreenResult {
	if result := c.localUpdate(state, event); result != nil {
		return *result
	}

	result := c.screen.Update(state, event)
	if result.Screen != nil {
		newBack := NewHistory(*result.Screen)
		newBack.history = &c.screen
		newScreen := newBack.ToScreen()
		result.Screen = &newScreen
	}
	return result
}

func (c *History) localUpdate(state state.UIState, event core.ScreenEvent) *core.ScreenResult {
	if event.Key == "b" && c.history != nil {
		newBack := NewHistory(*c.history)
		newScreen := newBack.ToScreen()
		result := core.NewScreenResult(state, &newScreen)
		return &result
	}

	return nil
}

func (c *History) view(state state.UIState) core.ViewModel {
	vm := c.screen.View(state)
	if c.history != nil {
		page := fmt.Sprintf("back: %s", c.history.Name())
		footer := core.NewLines(
			core.LineJump(),
			core.NewLine(page, core.ModePadding(core.Right)),
		)
		vm.Footer = append(vm.Footer, footer...)
	}
	return vm
}
