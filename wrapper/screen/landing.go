package wrapper_screen

import (
	"fmt"
	"unicode/utf8"

	"github.com/Rafael24595/go-terminal/engine/core/input"
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

	options := input.NewMenuOptions(
		input.NewMenuOption(text.NewFragment("Option Article"), NewTestArticle),
		input.NewMenuOption(text.NewFragment("Option TextArea"), NewTestTextArea),
		input.NewMenuOption(text.NewFragment("Option Table"), NewTestTable),
		input.NewMenuOption(text.NewFragment("Option Modal"), NewTestModal),
		input.NewMenuOption(text.NewFragment("Option Check"), NewTestCheck),
	)

	optsSize := len(options)

	for i := range 30 {
		options = append(options,
			input.NewMenuOption(text.NewFragment(fmt.Sprintf("Option %d", i+1+optsSize)), NewTestTextArea),
		)
	}

	return primitive.NewIndexMenu().
		SetName("menu - tortor").
		SetMeta(marker.NumericIndex).
		AddTitle(title...).
		AddOptions(options...).
		SetCursor(0).
		ToScreen()
}
