package marker

var BracketsCheck = CheckMeta{
	Open:      "[",
	Close:     "]",
	Checked:   "x",
	Unchecked: " ",
}

type CheckMeta struct {
	Open      string
	Close     string
	Checked   string
	Unchecked string
}

