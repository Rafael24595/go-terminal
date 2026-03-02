package wrapper_screen

import (
	drawable_table "github.com/Rafael24595/go-terminal/engine/core/drawable/table"
	"github.com/Rafael24595/go-terminal/engine/core/screen"
	"github.com/Rafael24595/go-terminal/engine/core/screen/commons"
	"github.com/Rafael24595/go-terminal/engine/core/style"
	"github.com/Rafael24595/go-terminal/engine/core/table"
	"github.com/Rafael24595/go-terminal/engine/core/text"
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
	return commons.NewTable[Language]().
		SetName("article - ipsum").
		DefinePadding(drawable_table.Center).
		EnableAction(func(drawable_table.Cursor) {}).
		AddTitle(
			text.LineFromString("Donec massa sem"),
			text.NewLine("=", style.SpecFromKind(style.SpcKindFill)),
			text.LineJump(),
		).
		DefineHeaders(headers...).
		AddItems(parser, rows...).
		ToScreen()
}
