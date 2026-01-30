package wrapper_screen

import (
	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/core/screen"
	"github.com/Rafael24595/go-terminal/engine/core/screen/commons"
)

func NewLanding() screen.Screen {
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

	options := commons.NewMenuOptions(
		commons.NewMenuOption(core.LineFromString("Option 0"), NewTestArticle),
		commons.NewMenuOption(core.LineFromString("Option 1"), NewTestTextArea),
		commons.NewMenuOption(core.LineFromString("Option 2"), NewTestArticle),
		commons.NewMenuOption(core.LineFromString("Option 3"), NewTestArticle),
	)

	return commons.NewIndexMenu().
		SetName("menu - tortor").
		AddTitle(title...).
		AddOptions(options...).
		SetCursor(0).
		ToScreen()
}
