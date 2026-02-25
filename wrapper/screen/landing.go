package wrapper_screen

import (
	"fmt"
	"unicode/utf8"

	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/core/screen"
	"github.com/Rafael24595/go-terminal/engine/core/screen/commons"
	"github.com/Rafael24595/go-terminal/engine/core/style"
)

func NewLanding() screen.Screen {
	textTitle := "Sed facilisis, leo sit amet molestie congue, justo risus bibendum tortor"
	sizeTitle := utf8.RuneCountInString(textTitle)
	

	title := core.NewLines(
		core.NewLine(
			textTitle,
			style.SpecFromKind(style.SpcKindPaddingRight),
		),
		core.NewLine(
			"-",
			style.SpecFill(uint(sizeTitle)),
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
