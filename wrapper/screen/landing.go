package wrapper_screen

import "github.com/Rafael24595/go-terminal/engine/core"

func NewLanding() core.Screen {
	title := core.NewLines(
		core.NewLine(
			"Sed facilisis, leo sit amet molestie congue, justo risus bibendum tortor",
			core.ModePadding(core.Right),
		),
		core.NewLine(
			"-",
			core.ModePadding(core.FillUp),
		),
	)

	options := NewMenuOptions(
		NewMenuOption(core.BasicLine("Option 0")),
		NewMenuOption(core.BasicLine("Option 1")),
		NewMenuOption(core.BasicLine("Option 2")),
		NewMenuOption(core.BasicLine("Option 3")),
	)

	return NewIndexMenu().
		AddTitle(title...).
		AddOptions(options...).
		SetCursor(0).
		ToScreen()
}
