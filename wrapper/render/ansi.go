package wrapper_render

const (
	Reset = "\033[0m"

	// Styles
	Bold      = "\033[1m"
	Dim       = "\033[2m"
	Italic    = "\033[3m"
	Underline = "\033[4m"
	Blink     = "\033[5m"
	Reverse   = "\033[7m"

	// Disable styles
	NoBold      = "\033[22m"
	NoItalic    = "\033[23m"
	NoUnderline = "\033[24m"
	NoBlink     = "\033[25m"
	NoReverse   = "\033[27m"
)
