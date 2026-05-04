package pass

import (
	"fmt"
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"

	screen_test "github.com/Rafael24595/go-reacterm-core/test/engine/app/screen"
)

func TestValidateStructure_ValidNode(t *testing.T) {
	node := screen_test.MockScreen{
		Name:       "home",
		Definition: &screen.Definition{},
		Update: func(*state.UIState, screen.ScreenEvent) screen.Result {
			return screen.Result{}
		},
		View: func(state.UIState) viewmodel.ViewModel {
			return viewmodel.ViewModel{}
		},
	}.ToNode()

	pass := ValidateStructure()
	_, err := pass(node)

	assert.Nil(t, err)
}

func TestValidateStructure_EmptyName(t *testing.T) {
	node := screen.Node{
		Screen: screen.Screen{
			Name: "",
			Definition: func() screen.Definition {
				return screen.Definition{}
			},
			Update: func(*state.UIState, screen.ScreenEvent) screen.Result {
				return screen.Result{}
			},
			View: func(state.UIState) viewmodel.ViewModel {
				return viewmodel.ViewModel{}
			},
		},
	}

	pass := ValidateStructure()
	_, err := pass(node)

	assert.Error(t, err)
	assert.Equal(t, err_name, err.Error())
}

func TestValidateStructure_NilDefinition(t *testing.T) {
	name := "home"

	node := screen.Node{
		Screen: screen.Screen{
			Name: name,
			Update: func(*state.UIState, screen.ScreenEvent) screen.Result {
				return screen.Result{}
			},
			View: func(state.UIState) viewmodel.ViewModel {
				return viewmodel.ViewModel{}
			},
		},
	}

	pass := ValidateStructure()
	_, err := pass(node)

	assert.Error(t, err)
	assert.Equal(t, fmt.Sprintf(errf_definition, name), err.Error())
}

func TestValidateStructure_NilUpdate(t *testing.T) {
	name := "home"

	node := screen.Node{
		Screen: screen.Screen{
			Name: name,
			Definition: func() screen.Definition {
				return screen.Definition{}
			},
			View: func(state.UIState) viewmodel.ViewModel {
				return viewmodel.ViewModel{}
			},
		},
	}

	pass := ValidateStructure()
	_, err := pass(node)

	assert.Error(t, err)
	assert.Equal(t, fmt.Sprintf(errf_update, name), err.Error())
}

func TestValidateStructure_NilView(t *testing.T) {
	name := "home"

	node := screen.Node{
		Screen: screen.Screen{
			Name: name,
			Definition: func() screen.Definition {
				return screen.Definition{}
			},
			Update: func(*state.UIState, screen.ScreenEvent) screen.Result {
				return screen.Result{}
			},
		},
	}

	pass := ValidateStructure()
	_, err := pass(node)

	assert.Error(t, err)
	assert.Equal(t, fmt.Sprintf(errf_view, name), err.Error())
}
