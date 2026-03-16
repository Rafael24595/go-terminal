package primitive

import (
	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/core/drawable/action"
	"github.com/Rafael24595/go-terminal/engine/core/drawable/stack"
	"github.com/Rafael24595/go-terminal/engine/core/screen"
)

type MapScreen struct {
	screen  screen.Screen
	actions []action.Action
}

func NewMapScreen(screen screen.Screen) *MapScreen {
	return &MapScreen{
		screen:  screen,
		actions: make([]action.Action, 0),
	}
}

func (c *MapScreen) PushAction(actions ...action.Action) *MapScreen {
	for _, a := range actions {
		c.actions = append(c.actions, a)
	}
	return c
}

func (c *MapScreen) ToScreen() screen.Screen {
	return screen.Screen{
		Name:       c.screen.Name,
		Definition: c.screen.Definition,
		Update:     c.update,
		View:       c.view,
		Stack:      c.screen.Stack,
	}
}

func (c *MapScreen) update(state *state.UIState, event screen.ScreenEvent) screen.ScreenResult {
	result := c.screen.Update(state, event)
	if result.Screen != nil {
		newScreen := NewMapScreen(*result.Screen).
			PushAction(c.actions...).
			ToScreen()
		result.Screen = &newScreen
	}
	return result
}

func (c *MapScreen) view(state state.UIState) core.ViewModel {
	vm := c.screen.View(state)

	header := vm.Header.Items()
	lines := vm.Lines.Items()
	footer := vm.Footer.Items()

	for _, a := range c.actions {
		if a.Focus == action.FocusNone {
			continue
		}

		if a.Focus.HasAny(action.FocusHeader) {
			header = action.ApplyAction(a, header...)
		}

		if a.Focus.HasAny(action.FocusBody) {
			lines = action.ApplyAction(a, lines...)
		}

		if a.Focus.HasAny(action.FocusFooter) {
			footer = action.ApplyAction(a, footer...)
		}
	}

	vm.Header = stack.NewStackDrawable().Shift(header...)
	vm.Lines = stack.NewStackDrawable().Shift(lines...)
	vm.Footer = stack.NewStackDrawable().Shift(footer...)

	return vm
}
