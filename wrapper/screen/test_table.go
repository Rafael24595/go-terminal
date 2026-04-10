package wrapper_screen

import (
	"github.com/Rafael24595/go-terminal/engine/app/screen"
	"github.com/Rafael24595/go-terminal/engine/app/screen/primitive"
	"github.com/Rafael24595/go-terminal/engine/model/table"
	"github.com/Rafael24595/go-terminal/engine/render/style"
	"github.com/Rafael24595/go-terminal/engine/render/text"
)

type Language struct {
	Name       string
	Creator    string
	Year       int
	Version    string
	Popularity float64
	Backend    bool
}

var headers = table.StructHeaders[Language]()

var rows = []Language{
	{
		Name:       "Go",
		Creator:    "Google",
		Year:       2009,
		Version:    "1.21",
		Popularity: 9.8,
		Backend:    true,
	},
	{
		Name:       "TypeScript",
		Creator:    "Microsoft",
		Year:       2012,
		Version:    "5.3",
		Popularity: 9.2,
		Backend:    false,
	},
	{
		Name:       "Rust",
		Creator:    "Mozilla",
		Year:       2010,
		Version:    "1.75",
		Popularity: 8.9,
		Backend:    true,
	},
	{
		Name:       "Python",
		Creator:    "Guido van Rossum",
		Year:       1991,
		Version:    "3.12",
		Popularity: 9.9,
		Backend:    true,
	},
	{
		Name:       "Swift",
		Creator:    "Apple",
		Year:       2014,
		Version:    "5.9",
		Popularity: 7.5,
		Backend:    false,
	},
}

func parser(lang Language) []table.Field {
	return table.StructFieds(lang)
}

func NewTestTable() screen.Screen {
	return primitive.NewTable[Language]().
		SetName("article - ipsum").
		DefinePadding(style.Center).
		EnableAction().
		AddTitle(
			*text.NewLine("Donec massa sem"),
			*text.NewLine("=", style.SpecFromKind(style.SpcKindFill)),
		).
		DefineHeaders(headers...).
		AddItems(parser, rows...).
		ToScreen()
}
