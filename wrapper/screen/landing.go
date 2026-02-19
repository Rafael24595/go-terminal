package wrapper_screen

import (
	"fmt"

	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/core/screen"
	"github.com/Rafael24595/go-terminal/engine/core/screen/commons"
	"github.com/Rafael24595/go-terminal/engine/core/style"
)

func NewLanding() screen.Screen {
	title := core.NewLines(
		core.NewLine(
			"Sed facilisis, leo sit amet molestie congue, justo risus bibendum tortor",
			style.SpecFromKind(style.SpcKindPaddingRight),
		),
		core.NewLine(
			"-",
			style.SpecFromKind(style.SpcKindFillUp),
		),
	)

	options := commons.NewMenuOptions(
		commons.NewMenuOption(core.LineFromString("Option Article"), NewTestArticle),
		commons.NewMenuOption(core.LineFromString("Option TextArea"), NewTestTextArea),
		commons.NewMenuOption(core.LineFromString("Option Table"), NewTestTable),
	)

	for i := range 30 {
		options = append(options,
			commons.NewMenuOption(core.LineFromString(fmt.Sprintf("Option %d", i+1)), NewTestTextArea),
		)
	}

	return commons.NewIndexMenu().
		SetName("menu - tortor").
		SetIndex(commons.NumericIndex).
		AddTitle(title...).
		AddOptions(options...).
		SetCursor(0).
		ToScreen()
}
