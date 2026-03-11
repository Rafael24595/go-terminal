package screen_test

import (
	"testing"

	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/core/screen"
	"github.com/Rafael24595/go-terminal/test/support/assert"
)

type MockScreen struct {
	Name       string
	Definition func() screen.Definition
	Update     func(*state.UIState, screen.ScreenEvent) screen.ScreenResult
	View       func(state.UIState) core.ViewModel
}

func (t MockScreen) ToScreen() screen.Screen {
	return screen.Screen{
		Name: func() string {
			return t.Name
		},
		Definition: func() screen.Definition {
			if t.Definition != nil {
				return t.Definition()
			}

			return screen.DefinitionFromKeys()
		},
		Update: func(s *state.UIState, e screen.ScreenEvent) screen.ScreenResult {
			if t.Update != nil {
				return t.Update(s, e)
			}

			return screen.ScreenResultFromUIState(s)
		},
		View: func(s state.UIState) core.ViewModel {
			if t.View != nil {
				return t.View(s)
			}

			return *core.ViewModelFromUIState(s)
		},
	}
}

func Helper_ToScreen(t *testing.T, screen screen.Screen) {
	t.Helper()

	assert.NotNil(t, screen.Name, "Screen.Name")
	assert.NotNil(t, screen.View, "Screen.View should be set")
	assert.NotNil(t, screen.Update, "Screen.Update should be set")
}
