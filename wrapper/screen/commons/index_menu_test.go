package wrapper_commons

import (
	"testing"

	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/test/support/assert"
	wrapper_terminal "github.com/Rafael24595/go-terminal/wrapper/terminal"
)

func TestIndexMenu_ToScreen(t *testing.T) {
	menu := NewIndexMenu().
		SetName("base").
		AddTitle(core.LineFromString("Welcome")).
		AddOptions(
			NewMenuOption(
				core.LineFromString("Option 1"),
				func() core.Screen { return core.Screen{} },
			),
		)

	screen := menu.ToScreen()

	Helper_ToScreen(t, screen)

	assert.Equal(t, screen.Name(), "base")
}

func TestIndexMenu_DefaultValues(t *testing.T) {
	menu := NewIndexMenu()

	assert.Equal(t, menu.reference, default_index_menu_name)
	assert.Equal(t, len(menu.title), 0)
	assert.Equal(t, len(menu.options), 0)
	assert.Equal(t, menu.cursor, uint(0))
}

func TestIndexMenu_AddTitleAndOptions(t *testing.T) {
	menu := NewIndexMenu().
		AddTitle(core.LineFromString("Title 1")).
		AddOptions(
			NewMenuOption(
				core.LineFromString("Option 1"),
				func() core.Screen { return core.Screen{} },
			),
			NewMenuOption(
				core.LineFromString("Option 2"),
				func() core.Screen { return core.Screen{} },
			),
		)

	assert.Equal(t, len(menu.title), 1)
	assert.Equal(t, len(menu.options), 2)
}

func TestIndexMenu_SetCursor_Clamp(t *testing.T) {
	menu := NewIndexMenu().
		AddOptions(
			NewMenuOption(
				core.LineFromString("A"),
				func() core.Screen { return core.Screen{} },
			),
			NewMenuOption(
				core.LineFromString("B"),
				func() core.Screen { return core.Screen{} },
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
			NewMenuOption(
				core.LineFromString("A"),
				func() core.Screen { return core.Screen{} },
			),
			NewMenuOption(
				core.LineFromString("B"),
				func() core.Screen { return core.Screen{} },
			),
		)

	screen := menu.ToScreen()

	assert.Equal(t, menu.cursor, uint(0))

	screen.Update(
		*state.NewUIState(),
		core.ScreenEvent{Key: wrapper_terminal.ARROW_DOWN},
	)
	assert.Equal(t, menu.cursor, uint(1))

	screen.Update(
		*state.NewUIState(),
		core.ScreenEvent{Key: wrapper_terminal.ARROW_UP},
	)
	assert.Equal(t, menu.cursor, uint(0))
}

func TestIndexMenu_Action(t *testing.T) {
	expected := core.Screen{
		Name: func() string { return "next" },
	}

	menu := NewIndexMenu().
		AddOptions(
			NewMenuOption(
				core.LineFromString("Go"),
				func() core.Screen { return expected },
			),
		)

	screen := menu.ToScreen()
	result := screen.Update(
		*state.NewUIState(),
		core.ScreenEvent{Key: "\n"},
	)

	assert.NotNil(t, result.Screen)
	assert.Equal(t, result.Screen.Name(), "next")
}

func TestIndexMenu_ViewCursor(t *testing.T) {
	menu := NewIndexMenu().
		AddOptions(
			NewMenuOption(
				core.LineFromString("A"),
				func() core.Screen { return core.Screen{} },
			),
			NewMenuOption(
				core.LineFromString("B"),
				func() core.Screen { return core.Screen{} },
			),
		)

	menu.cursor = 1
	vm := menu.view(state.UIState{})

	assert.NotNil(t, vm.Cursor)
	assert.Equal(t, *vm.Cursor, uint(1))

	assert.Equal(t, vm.Lines[1].Text[0].Text, ">")
}
