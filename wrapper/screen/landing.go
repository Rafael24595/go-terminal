package wrapper_screen

import (
	"fmt"
	"unicode/utf8"

	"github.com/Rafael24595/go-terminal/engine/core/marker"
	"github.com/Rafael24595/go-terminal/engine/core/screen"
	"github.com/Rafael24595/go-terminal/engine/core/screen/primitive"
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

	options := primitive.NewMenuOptions(
		primitive.NewMenuOption(text.NewFragment("Option Article"), NewTestArticle),
		primitive.NewMenuOption(text.NewFragment("Option TextArea"), NewTestTextArea),
		primitive.NewMenuOption(text.NewFragment("Option Table"), NewTestTable),
		primitive.NewMenuOption(text.NewFragment("Option Modal"), NewTestModal),
	)

	optsSize := len(options)

	for i := range 30 {
		options = append(options,
			primitive.NewMenuOption(text.NewFragment(fmt.Sprintf("Option %d", i+1+optsSize)), NewTestTextArea),
		)
	}

	return primitive.NewIndexMenu().
		SetName("menu - tortor").
		SetIndex(marker.NumericIndex).
		AddTitle(title...).
		AddOptions(options...).
		SetCursor(0).
		ToScreen()
}
