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
		Name: "home",
		Keys: &screen.Definition{},
		Tick: func(*state.UIState, screen.Event) screen.Result {
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
			Keys: func() screen.Definition {
				return screen.Definition{}
			},
			Tick: func(*state.UIState, screen.Event) screen.Result {
				return screen.Result{}
			},
			View: func(state.UIState) viewmodel.ViewModel {
				return viewmodel.ViewModel{}
			},
		},
	}

	pass := ValidateStructure()
	_, err := pass(node)

	assert.NotNil(t, err)
	assert.Equal(t, err_name, err.Error())
}

func TestValidateStructure_NilKeys(t *testing.T) {
	name := "home"

	node := screen.Node{
		Name: name,
		Screen: screen.Screen{
			Tick: func(*state.UIState, screen.Event) screen.Result {
				return screen.Result{}
			},
			View: func(state.UIState) viewmodel.ViewModel {
				return viewmodel.ViewModel{}
			},
		},
	}

	pass := ValidateStructure()
	_, err := pass(node)

	assert.NotNil(t, err)
	assert.Equal(t, fmt.Sprintf(errf_keys, name), err.Error())
}

func TestValidateStructure_NilTick(t *testing.T) {
	name := "home"

	node := screen.Node{
		Name: name,
		Screen: screen.Screen{
			Keys: func() screen.Definition {
				return screen.Definition{}
			},
			View: func(state.UIState) viewmodel.ViewModel {
				return viewmodel.ViewModel{}
			},
		},
	}

	pass := ValidateStructure()
	_, err := pass(node)

	assert.NotNil(t, err)
	assert.Equal(t, fmt.Sprintf(errf_tick, name), err.Error())
}

func TestValidateStructure_NilView(t *testing.T) {
	name := "home"

	node := screen.Node{
		Name: name,
		Screen: screen.Screen{
			Keys: func() screen.Definition {
				return screen.Definition{}
			},
			Tick: func(*state.UIState, screen.Event) screen.Result {
				return screen.Result{}
			},
		},
	}

	pass := ValidateStructure()
	_, err := pass(node)

	assert.NotNil(t, err)
	assert.Equal(t, fmt.Sprintf(errf_view, name), err.Error())
}
