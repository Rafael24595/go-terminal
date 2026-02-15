package commons

import (
	"testing"

	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/core/key"
	"github.com/Rafael24595/go-terminal/engine/core/screen"
	"github.com/Rafael24595/go-terminal/engine/terminal"
	"github.com/Rafael24595/go-terminal/test/support/assert"
)

func TestIndexMenu_ToScreen(t *testing.T) {
	menu := NewIndexMenu().
		SetName("base").
		AddTitle(core.LineFromString("Welcome")).
		AddOptions(
			NewMenuOption(
				core.LineFromString("Option 1"),
				func() screen.Screen { return screen.Screen{} },
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
				func() screen.Screen { return screen.Screen{} },
			),
			NewMenuOption(
				core.LineFromString("Option 2"),
				func() screen.Screen { return screen.Screen{} },
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
				func() screen.Screen { return screen.Screen{} },
			),
			NewMenuOption(
				core.LineFromString("B"),
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
			NewMenuOption(
				core.LineFromString("A"),
				func() screen.Screen { return screen.Screen{} },
			),
			NewMenuOption(
				core.LineFromString("B"),
				func() screen.Screen { return screen.Screen{} },
			),
		)

	scrn := menu.ToScreen()

	assert.Equal(t, menu.cursor, uint(0))

	scrn.Update(
		*state.NewUIState(),
		screen.ScreenEvent{Key: *key.NewKeyCode(key.ActionArrowDown)},
	)
	assert.Equal(t, menu.cursor, uint(1))

	scrn.Update(
		*state.NewUIState(),
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
			NewMenuOption(
				core.LineFromString("Go"),
				func() screen.Screen { return expected },
			),
		)

	scrn := menu.ToScreen()
	result := scrn.Update(
		*state.NewUIState(),
		screen.ScreenEvent{Key: *key.NewKeyCode(key.ActionEnter)},
	)

	assert.NotNil(t, result.Screen)
	assert.Equal(t, result.Screen.Name(), "next")
}

func TestIndexMenu_ViewCursor(t *testing.T) {
	menu := NewIndexMenu().
		AddOptions(
			NewMenuOption(
				core.LineFromString("A"),
				func() screen.Screen { return screen.Screen{} },
			),
			NewMenuOption(
				core.LineFromString("B"),
				func() screen.Screen { return screen.Screen{} },
			),
		)

	menu.cursor = 1
	vm := menu.view(state.UIState{})

	vm.Lines.Init(terminal.Winsize{})
	lines, _ := vm.Lines.Draw()

	assert.NotNil(t, vm.Cursor)
	assert.True(t, vm.Cursor.Enabled)
	assert.Equal(t, vm.Cursor.Cursor, uint(1))

	assert.Equal(t, lines[1].Text[0].Text, ">")
}
