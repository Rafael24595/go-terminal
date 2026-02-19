package wrapper_screen

import (
	"reflect"

	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/core/screen"
	"github.com/Rafael24595/go-terminal/engine/core/screen/commons"
	"github.com/Rafael24595/go-terminal/engine/core/style"
)

type Language struct {
	Name       string
	Creator    string
	Year       int
	Version    string
	Popularity float64
	Backend    bool
}

var headers = headersStruct[Language]()

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

func parser(lang Language) []commons.Field {
	return parseStruct(lang)
}

func NewTestTable() screen.Screen {
	return commons.NewTable[Language]().
		SetName("article - ipsum").
		AddTitle(
			core.LineFromString("Donec massa sem"),
			core.NewLine("=", style.SpecFromKind(style.SpcKindFill)),
			core.LineJump(),
		).
		DefineHeaders(headers...).
		AddItems(parser, rows...).
		ToScreen()
}

func headersStruct[T any]() []string {
	var zero T

	headers := make([]string, 0)
	for _, v := range parseStruct(zero) {
		headers = append(headers, v.Header)
	}
	return headers
}

func parseStruct(s any) []commons.Field {
	v := reflect.ValueOf(s)
	t := reflect.TypeOf(s)

	if t.Kind() == reflect.Pointer {
		v = v.Elem()
		t = t.Elem()
	}

	var result []commons.Field

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i).Interface()

		result = append(result, commons.Field{
			Header: field.Name,
			Value:  value,
		})
	}

	return result
}
