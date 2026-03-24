package screen_test

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	
	"github.com/Rafael24595/go-terminal/engine/app/screen"
	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/app/viewmodel"
	"github.com/Rafael24595/go-terminal/engine/commons/structure/set"
)

type MockScreen struct {
	Name       string
	Definition *screen.Definition
	Update     func(*state.UIState, screen.ScreenEvent) screen.ScreenResult
	View       func(state.UIState) viewmodel.ViewModel
	Stack      set.Set[string]
}

func (t MockScreen) ToScreen() screen.Screen {
	return screen.Screen{
		Name: func() string {
			return t.Name
		},
		Definition: func() screen.Definition {
			if t.Definition != nil {
				return *t.Definition
			}

			return screen.DefinitionFromKeys()
		},
		Update: func(s *state.UIState, e screen.ScreenEvent) screen.ScreenResult {
			if t.Update != nil {
				return t.Update(s, e)
			}

			return screen.ScreenResultFromUIState(s)
		},
		View: func(s state.UIState) viewmodel.ViewModel {
			if t.View != nil {
				return t.View(s)
			}

			return *viewmodel.ViewModelFromUIState(s)
		},
		Stack: func() set.Set[string] {
			if t.Stack != nil {
				return t.Stack
			}

			return set.SetFrom(t.Name)
		},
	}
}

func Helper_ToScreen(t *testing.T, screen screen.Screen) {
	t.Helper()

	assert.NotNil(t, screen.Name, "Screen.Name")
	assert.NotNil(t, screen.View, "Screen.View should be set")
	assert.NotNil(t, screen.Update, "Screen.Update should be set")
	assert.NotNil(t, screen.Stack, "Screen.Stack should be set")
}
