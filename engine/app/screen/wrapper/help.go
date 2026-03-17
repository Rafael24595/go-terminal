package wrapper

import (
	"github.com/Rafael24595/go-terminal/engine/app/screen"
	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/app/viewmodel"
	"github.com/Rafael24595/go-terminal/engine/model/key"
)

type Help struct {
	helpMode bool
	screen   screen.Screen
}

func NewHelp(screen screen.Screen) *Help {
	return &Help{
		helpMode: false,
		screen:   screen,
	}
}

func (c *Help) ToScreen() screen.Screen {
	return screen.Screen{
		Name:       c.screen.Name,
		Definition: c.screen.Definition,
		Update:     c.update,
		View:       c.view,
		Stack:      c.screen.Stack,
	}
}

func (c *Help) update(state *state.UIState, event screen.ScreenEvent) screen.ScreenResult {
	requiredKey := screen.IsKeyRequired(c.screen.Definition(), event.Key)

	if requiredKey {
		result := c.screen.Update(state, event)
		if result.Screen != nil {
			newWrapper := NewHelp(*result.Screen)
			newWrapper.helpMode = c.helpMode
			newScreen := newWrapper.ToScreen()
			result.Screen = &newScreen
		}

		c.helpMode = state.Helper.ShowHelp
		return result
	}

	if event.Key.Code == key.CustomActionHelp {
		c.helpMode = !c.helpMode
	}

	state.Helper.ShowHelp = c.helpMode
	return screen.ScreenResultFromUIState(state)
}

func (c *Help) view(state state.UIState) viewmodel.ViewModel {
	vm := c.screen.View(state)

	vm.Helper.Show = c.helpMode

	return vm
}
