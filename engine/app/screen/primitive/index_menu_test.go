package primitive

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-terminal/engine/app/pager"
	"github.com/Rafael24595/go-terminal/engine/app/screen"
	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/model/input"
	"github.com/Rafael24595/go-terminal/engine/model/key"
	"github.com/Rafael24595/go-terminal/engine/render/text"
	"github.com/Rafael24595/go-terminal/engine/terminal"

	screen_test "github.com/Rafael24595/go-terminal/test/engine/app/screen"
)

func TestIndexMenu_ToScreen(t *testing.T) {
	menu := NewIndexMenu().
		SetName("base").
		AddTitle(
			text.NewLine("Welcome"),
		).
		AddOptions(
			input.NewMenuOption(
				"opt_1",
				text.NewFragment("Option 1"),
				func() screen.Screen { return screen.Screen{} },
			),
		)

	screen := menu.ToScreen()

	screen_test.Helper_ToScreen(t, screen)

	assert.Equal(t, screen.Name(), "base")
}

func TestIndexMenu_Stack(t *testing.T) {
	stack := NewIndexMenu().
		ToScreen().
		Stack()

	assert.True(t, stack.Has(default_index_menu_name))
}

func TestIndexMenu_DefaultValues(t *testing.T) {
	menu := NewIndexMenu()

	assert.Equal(t, menu.reference, default_index_menu_name)
	assert.Len(t, 0, menu.title)
	assert.Len(t, 0, menu.options)
	assert.Equal(t, menu.cursor, uint(0))
}

func TestIndexMenu_AddTitleAndOptions(t *testing.T) {
	menu := NewIndexMenu().
		AddTitle(
			text.NewLine("Title 1"),
		).
		AddOptions(
			input.NewMenuOption(
				"opt_1",
				text.NewFragment("Option 1"),
				func() screen.Screen { return screen.Screen{} },
			),
			input.NewMenuOption(
				"opt_2",
				text.NewFragment("Option 2"),
				func() screen.Screen { return screen.Screen{} },
			),
		)

	assert.Len(t, 1, menu.title)
	assert.Len(t, 2, menu.options)
}

func TestIndexMenu_SetCursor_Clamp(t *testing.T) {
	menu := NewIndexMenu().
		AddOptions(
			input.NewMenuOption(
				"opt_a",
				text.NewFragment("A"),
				func() screen.Screen { return screen.Screen{} },
			),
			input.NewMenuOption(
				"opt_b",
				text.NewFragment("B"),
				func() screen.Screen { return screen.Screen{} },
			),
		)

	menu.SetCursor(10)
	assert.Equal(t, menu.cursor, uint(1))

	menu.SetCursor(0)
	assert.Equal(t, menu.cursor, uint(0))
}

func TestIndexMenu_SetCursor_Empty(t *testing.T) {
	menu := NewIndexMenu()
	menu.SetCursor(5)

	assert.Equal(t, menu.cursor, uint(0))
}

func TestIndexMenu_CursorNavigation(t *testing.T) {
	menu := NewIndexMenu().
		AddOptions(
			input.NewMenuOption(
				"opt_a",
				text.NewFragment("A"),
				func() screen.Screen { return screen.Screen{} },
			),
			input.NewMenuOption(
				"opt_b",
				text.NewFragment("B"),
				func() screen.Screen { return screen.Screen{} },
			),
		)

	scrn := menu.ToScreen()

	assert.Equal(t, menu.cursor, uint(0))

	scrn.Update(
		state.NewUIState(),
		screen.ScreenEvent{Key: *key.NewKeyCode(key.ActionArrowDown)},
	)
	assert.Equal(t, menu.cursor, uint(1))

	scrn.Update(
		state.NewUIState(),
		screen.ScreenEvent{Key: *key.NewKeyCode(key.ActionArrowUp)},
	)
	assert.Equal(t, menu.cursor, uint(0))
}

func TestIndexMenu_Action(t *testing.T) {
	expected := screen.Screen{
		Name: func() string { return "next" },
	}

	menu := NewIndexMenu().
		AddOptions(
			input.NewMenuOption(
				"opt_go",
				text.NewFragment("Go"),
				func() screen.Screen { return expected },
			),
		)

	scrn := menu.ToScreen()
	result := scrn.Update(
		state.NewUIState(),
		screen.ScreenEvent{Key: *key.NewKeyCode(key.ActionEnter)},
	)

	assert.NotNil(t, result.Screen)
	assert.Equal(t, result.Screen.Name(), "next")
}

func TestIndexMenu_ViewCursor(t *testing.T) {
	menu := NewIndexMenu().
		AddOptions(
			input.NewMenuOption(
				"opt_a",
				text.NewFragment("A"),
				func() screen.Screen { return screen.Screen{} },
			),
			input.NewMenuOption(
				"opt_b",
				text.NewFragment("B"),
				func() screen.Screen { return screen.Screen{} },
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
	lines, _ := kernel.Draw(terminal.Winsize{Cols: 10, Rows: 2})

	assert.NotNil(t, vm.Pager)
	assert.Equal(t, pager.CodePredicateFocus, vm.Pager.Predicate.Code)
	assert.True(t, vm.Pager.Predicate.Func(*stt, ctx))

	assert.Equal(t, "- A", text.LineToString(lines[0]))
	assert.Equal(t, "> B", text.LineToString(lines[1]))
}
