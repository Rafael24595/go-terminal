package wrapper_screen

import (
	"fmt"

	"github.com/Rafael24595/go-terminal/engine/app/screen"
	"github.com/Rafael24595/go-terminal/engine/app/screen/primitive"
	"github.com/Rafael24595/go-terminal/engine/helper/runes"
	"github.com/Rafael24595/go-terminal/engine/model/input"
	"github.com/Rafael24595/go-terminal/engine/render/marker"
	"github.com/Rafael24595/go-terminal/engine/render/style"
	"github.com/Rafael24595/go-terminal/engine/render/text"
)

func NewLanding() screen.Screen {
	textTitle := "Sed facilisis, leo sit amet molestie congue, justo risus bibendum tortor"
	sizeTitle := runes.Measure(textTitle)

	title := []text.Line{
		*text.NewLine(
			textTitle,
			style.SpecFromKind(style.SpcKindPaddingRight),
		),
		*text.NewLine(
			"-",
			style.SpecFill(uint(sizeTitle)),
		),
		*text.EmptyLine(),
	}

	options := input.NewMenuOptions(
		input.NewMenuOption("opt_art", *text.NewFragment("Option Article"), NewTestArticle),
		input.NewMenuOption("opt_txt", *text.NewFragment("Option TextArea"), NewTestTextArea),
		input.NewMenuOption("opt_tbl", *text.NewFragment("Option Table"), NewTestTable),
		input.NewMenuOption("opt_mdl", *text.NewFragment("Option Modal"), NewTestModal),
		input.NewMenuOption("opt_chk", *text.NewFragment("Option Check"), NewTestCheck),
		input.NewMenuOption("opt_chk", *text.NewFragment("Option TextInput"), NewTestTextInput),
		input.NewMenuOption("opt_hsk", *text.NewFragment("Option HStack"), NewTestHStack),
	)

	optsSize := len(options)

	for i := range 30 {
		options = append(options,
			input.NewMenuOption(
				fmt.Sprintf("opt_%d", i),
				*text.NewFragment(fmt.Sprintf("Option %d", i+1+optsSize)),
				NewTestTextArea,
			),
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
