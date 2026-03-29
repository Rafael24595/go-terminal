package marker

type BoxSeparatorMeta struct {
	Top    string
	Bottom string
	Left   string
	Right  string
	//TODO: Deprecate?
	Space  string
}

var DefaultBoxSeparator = BoxSeparatorMeta{
	Top:    "-",
	Bottom: "-",
	Left:   "|",
	Right:  "|",
	Space:  " ",
}
