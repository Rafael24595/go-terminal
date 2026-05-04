package wrapper_screen

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/node/primitive/checkmenu"
	"github.com/Rafael24595/go-reacterm-core/engine/helper/runes"
	"github.com/Rafael24595/go-reacterm-core/engine/model/input"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

func NewTestCheck() screen.Node {
	textTitle := "Sed facilisis, leo sit amet molestie congue, justo risus bibendum tortor"
	sizeTitle := runes.Measure(textTitle)

	title := []text.Line{
		*text.NewLine(
			textTitle,
			style.SpecFromKind(style.SpcKindPaddingRight),
		),
		*text.NewLine(
			"-",
			style.SpecFill(sizeTitle),
		),
		*text.EmptyLine(),
	}

	options := []input.CheckOption{
		input.NewCheckOption("1", *text.NewFragment("Check 1")),
		input.NewCheckOption("2", *text.NewFragment("Check 2")),
		input.NewCheckOption("3", *text.NewFragment("Check 3")),
		input.NewCheckOption("4", *text.NewFragment("Check 4")),
	}

	return checkmenu.New().
		Name("menu - tortor").
		Limit(1).
		AddTitle(title...).
		AddOptions(options...).
		ToNode()
}
