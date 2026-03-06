package marker

type BoxSeparatorMeta struct {
	Top    string
	Bottom string
	Left   string
	Right  string
	Space  string
}

var DefaultBoxSeparator = BoxSeparatorMeta{
	Top:    "-",
	Bottom: "-",
	Left:   "|",
	Right:  "|",
	Space:  " ",
}
