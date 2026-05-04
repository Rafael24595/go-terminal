package indexmenu

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/app/pager"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/model/input"
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"

	screen_test "github.com/Rafael24595/go-reacterm-core/test/engine/app/screen"
)

func voidAction() screen.Node { return screen.Node{} }

func TestIndexMenu_ToNode(t *testing.T) {
	node := New().
		SetName("base").
		AddTitle(
			*text.NewLine("Welcome"),
		).
		AddOptions(
			input.NewMenuOption(
				"opt_1",
				*text.NewFragment("Option 1"),
				voidAction,
			),
		).
		ToNode()

	screen_test.Helper_ToNode(t, node)

	assert.Equal(t, node.Screen.Name, "base")
}

func TestIndexMenu_Stack(t *testing.T) {
	stack := New().ToNode().Stack

	assert.True(t, stack.Has(Name))
}

func TestIndexMenu_DefaultValues(t *testing.T) {
	menu := New()

	assert.Equal(t, menu.reference, Name)
	assert.Len(t, 0, menu.title)
	assert.Len(t, 0, menu.options)
	assert.Equal(t, menu.cursor, 0)
}

func TestIndexMenu_AddTitleAndOptions(t *testing.T) {
	menu := New().
		AddTitle(
			*text.NewLine("Title 1"),
		).
		AddOptions(
			input.NewMenuOption(
				"opt_1",
				*text.NewFragment("Option 1"),
				voidAction,
			),
			input.NewMenuOption(
				"opt_2",
				*text.NewFragment("Option 2"),
				voidAction,
			),
		)

	assert.Len(t, 1, menu.title)
	assert.Len(t, 2, menu.options)
}

func TestIndexMenu_SetCursor_Clamp(t *testing.T) {
	menu := New().
		AddOptions(
			input.NewMenuOption(
				"opt_a",
				*text.NewFragment("A"),
				voidAction,
			),
			input.NewMenuOption(
				"opt_b",
				*text.NewFragment("B"),
				voidAction,
			),
		)

	menu.SetCursor(10)
	assert.Equal(t, menu.cursor, 1)

	menu.SetCursor(0)
	assert.Equal(t, menu.cursor, 0)
}

func TestIndexMenu_SetCursor_Empty(t *testing.T) {
	menu := New()
	menu.SetCursor(5)

	assert.Equal(t, menu.cursor, 0)
}

func TestIndexMenu_CursorNavigation(t *testing.T) {
	menu := New().
		AddOptions(
			input.NewMenuOption(
				"opt_a",
				*text.NewFragment("A"),
				voidAction,
			),
			input.NewMenuOption(
				"opt_b",
				*text.NewFragment("B"),
				voidAction,
			),
		)

	node := menu.ToNode()

	assert.Equal(t, menu.cursor, 0)

	node.Screen.Update(
		state.NewUIState(),
		screen.ScreenEvent{Key: *key.NewKeyCode(key.ActionArrowDown)},
	)
	assert.Equal(t, menu.cursor, 1)

	node.Screen.Update(
		state.NewUIState(),
		screen.ScreenEvent{Key: *key.NewKeyCode(key.ActionArrowUp)},
	)
	assert.Equal(t, menu.cursor, 0)
}

func TestIndexMenu_Action(t *testing.T) {
	expected := screen.Node{
		Screen: screen.Screen{
			Name: "next",
		},
	}

	menu := New().
		AddOptions(
			input.NewMenuOption(
				"opt_go",
				*text.NewFragment("Go"),
				func() screen.Node { return expected },
			),
		)

	node := menu.ToNode()
	result := node.Screen.Update(
		state.NewUIState(),
		screen.ScreenEvent{Key: *key.NewKeyCode(key.ActionEnter)},
	)

	assert.NotNil(t, result.Node)
	assert.Equal(t, result.Node.Screen.Name, "next")
}

func TestIndexMenu_ViewCursor(t *testing.T) {
	menu := New().
		AddOptions(
			input.NewMenuOption(
				"opt_a",
				*text.NewFragment("A"),
				voidAction,
			),
			input.NewMenuOption(
				"opt_b",
				*text.NewFragment("B"),
				voidAction,
			),
		)

	stt := &state.UIState{}

	ctx := pager.PredicateContext{
		HasFocus: true,
	}

	menu.cursor = 1
	vm := menu.view(*stt)

	kernel := vm.Kernel.ToDrawable()

	kernel.Init()
	lines, _ := kernel.Draw(winsize.Winsize{Cols: 10, Rows: 2})

	assert.NotNil(t, vm.Pager)
	assert.Equal(t, pager.CodePredicateFocus, vm.Pager.Predicate.Code)
	assert.True(t, vm.Pager.Predicate.Func(*stt, ctx))

	assert.Equal(t, "- A", text.LineToString(&lines[0]))
	assert.Equal(t, "> B", text.LineToString(&lines[1]))
}
