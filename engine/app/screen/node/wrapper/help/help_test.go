package help

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"

	screen_test "github.com/Rafael24595/go-reacterm-core/test/engine/app/screen"
)

func TestHelp_ToNode(t *testing.T) {
	name := "base"
	mock := screen_test.MockScreen{
		Name: name,
	}

	node := New(mock.ToNode()).ToNode()
	screen_test.Helper_ToNode(t, node)

	assert.Equal(t, node.Screen.Name, name)
}

func TestHelp_Propagate(t *testing.T) {
	name := "base"
	mock := screen_test.MockScreen{
		Name: name,
	}

	node := New(mock.ToNode()).ToNode()
	screen_test.Helper_Propagate(t, name, 0, node)
}

func TestHelp_ToggleHelpKey(t *testing.T) {
	called := false

	mock := screen_test.MockScreen{}

	node := New(mock.ToNode()).ToNode()

	state := &state.UIState{}
	event := screen.ScreenEvent{
		Key: *key.NewKeyCode(key.CustomActionHelp),
	}

	node.Screen.Update(state, event)

	assert.True(t, state.Helper.ShowHelp)
	assert.False(t, called)
}

func TestHelp_DelegatesUpdateWhenKeyRequired(t *testing.T) {
	called := false

	ky := *key.NewKeyCode(key.CustomActionHelp)
	definition := screen.DefinitionFromKeys(ky)

	mock := screen_test.MockScreen{
		Definition: &definition,
		Update: func(s *state.UIState, e screen.ScreenEvent) screen.Result {
			called = true
			return screen.EmptyResult()
		},
	}

	node := New(mock.ToNode()).ToNode()

	state := &state.UIState{}
	event := screen.ScreenEvent{
		Key: ky,
	}

	node.Screen.Update(state, event)

	assert.False(t, state.Helper.ShowHelp)
	assert.True(t, called)
}

func TestHelp_WrapsReturnedScreen(t *testing.T) {
	called := false

	ky := *key.NewKeyCode(key.ActionEnter)
	definition := screen.DefinitionFromKeys(ky)

	mockNext := screen_test.MockScreen{
		Name: "next",
	}

	mockBase := screen_test.MockScreen{
		Definition: &definition,
		Update: func(s *state.UIState, _ screen.ScreenEvent) screen.Result {
			called = true
			next := mockNext.ToNode()
			return screen.Result{
				Node: &next,
			}
		},
	}

	help := New(mockBase.ToNode())
	wrapped := help.ToNode()

	stt := &state.UIState{}
	evt := screen.ScreenEvent{
		Key: ky,
	}

	wrapped.Screen.Update(stt, screen.ScreenEvent{
		Key: *key.NewKeyCode(key.CustomActionHelp),
	})

	assert.True(t, stt.Helper.ShowHelp)

	result := wrapped.Screen.Update(stt, evt)

	assert.True(t, called)
	assert.NotNil(t, result.Node)
	assert.Equal(t, "next", result.Node.Screen.Name)

	vm := result.Node.Screen.View(state.UIState{})

	assert.True(t, vm.Helper.Show)
}
