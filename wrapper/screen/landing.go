package wrapper_screen

import (
	"github.com/Rafael24595/go-terminal/engine/core"
	wrapper_commons "github.com/Rafael24595/go-terminal/wrapper/screen/commons"
)

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

	options := wrapper_commons.NewMenuOptions(
		wrapper_commons.NewMenuOption(core.LineFromString("Option 0"), NewTestArticle),
		wrapper_commons.NewMenuOption(core.LineFromString("Option 1"), NewTestArticle),
		wrapper_commons.NewMenuOption(core.LineFromString("Option 2"), NewTestArticle),
		wrapper_commons.NewMenuOption(core.LineFromString("Option 3"), NewTestArticle),
	)

	return wrapper_commons.NewIndexMenu().
		SetName("menu - tortor").
		AddTitle(title...).
		AddOptions(options...).
		SetCursor(0).
		ToScreen()
}
