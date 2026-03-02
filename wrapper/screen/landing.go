package wrapper_screen

import (
	"fmt"
	"unicode/utf8"

	"github.com/Rafael24595/go-terminal/engine/core/screen"
	"github.com/Rafael24595/go-terminal/engine/core/screen/commons"
	"github.com/Rafael24595/go-terminal/engine/core/style"
	"github.com/Rafael24595/go-terminal/engine/core/text"
)

func NewLanding() screen.Screen {
	textTitle := "Sed facilisis, leo sit amet molestie congue, justo risus bibendum tortor"
	sizeTitle := utf8.RuneCountInString(textTitle)

	title := text.NewLines(
		text.NewLine(
			textTitle,
			style.SpecFromKind(style.SpcKindPaddingRight),
		),
		text.NewLine(
			"-",
			style.SpecFill(uint(sizeTitle)),
		),
	)

	options := commons.NewMenuOptions(
		commons.NewMenuOption(text.LineFromString("Option Article"), NewTestArticle),
		commons.NewMenuOption(text.LineFromString("Option TextArea"), NewTestTextArea),
		commons.NewMenuOption(text.LineFromString("Option Table"), NewTestTable),
	)

	for i := range 30 {
		options = append(options,
			commons.NewMenuOption(text.LineFromString(fmt.Sprintf("Option %d", i+1)), NewTestTextArea),
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
