package history

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"

	screen_test "github.com/Rafael24595/go-reacterm-core/test/engine/app/screen"
)

func TestHistory_ToNode(t *testing.T) {
	name := "base"
	mock := screen_test.MockScreen{
		Name: name,
	}

	node := New(mock.ToNode()).ToNode()
	screen_test.Helper_ToNode(t, node)

	assert.Equal(t, node.Name, name)
}

func TestHistory_Propagate(t *testing.T) {
	name := "base"
	mock := screen_test.MockScreen{
		Name: name,
	}

	node := New(mock.ToNode()).ToNode()
	screen_test.Helper_Propagate(t, name, 0, node)
}
func TestHistory_BackNavigation(t *testing.T) {
	stt := &state.UIState{}

	mockBase := screen_test.MockScreen{
		Name: "base",
	}

	mockNext := screen_test.MockScreen{
		Name: "next",
		Update: func(s *state.UIState, e screen.Event) screen.Result {
			base := mockBase.ToNode()
			return screen.ResultFromNode(&base)
		},
	}

	node := New(mockNext.ToNode()).
		ToNode()

	assert.Equal(t, node.Name, "next")

	result := node.Screen.Update(stt, screen.Event{})
	assert.NotNil(t, result.Node)
	assert.Equal(t, result.Node.Name, "base")

	backResult := result.Node.Screen.Update(stt, screen.Event{
		Key: *key.NewKeyCode(key.CustomActionBack),
	})

	assert.NotNil(t, backResult.Node.Screen)
	assert.Equal(t, backResult.Node.Name, "next")
}

func TestHistory_ViewFooter(t *testing.T) {
	mock := screen_test.MockScreen{}
	node := mock.ToNode()

	h := New(node)

	vm := h.view(*state.NewUIState())

	footer := vm.Footer.ToUnit()
	footer.Drawable.Init()

	lines, _ := footer.Drawable.Draw(winsize.Winsize{})

	assert.Len(t, 0, lines)

	h.history = &node
	vm = h.view(*state.NewUIState())

	footer = vm.Footer.ToUnit()
	footer.Drawable.Init()

	lines, _ = footer.Drawable.Draw(winsize.Winsize{
		Rows: 3,
		Cols: 10,
	})

	assert.Len(t, 1, lines)
}
