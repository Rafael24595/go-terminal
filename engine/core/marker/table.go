package marker

type TableSeparatorMeta struct {
	Top    string
	Bottom string
	Center string
	Left   string
	Right  string
}

var DefaultTableSeparator = TableSeparatorMeta{
	Top:    "-",
	Bottom: "-",
	Center: " | ",
	Left:   "| ",
	Right:  " |",
}
