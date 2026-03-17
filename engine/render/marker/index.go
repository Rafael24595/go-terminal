package marker

type IndexKind int

const (
	Numeric IndexKind = iota
	Alphabetic
	Custom
)

var NumericIndex = KindIndex(Numeric)
var AlphabeticIndex = KindIndex(Alphabetic)

var GreaterIndex = CustomIndex(">", "-")
var HyphenIndex = CustomIndex("-", ">")

type IndexMeta struct {
	Kind   IndexKind
	Index  string
	Cursor string
}

func KindIndex(kind IndexKind) IndexMeta {
	return IndexMeta{
		Kind: kind,
	}
}

func CustomIndex(index string, cursor string) IndexMeta {
	return IndexMeta{
		Kind:   Custom,
		Index:  index,
		Cursor: cursor,
	}
}
