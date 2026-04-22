package wrapper_ansi

const (
	FullReset = ClearScreen + CursorHome
)

const (
	CursorHome  = "\x1b[H"
	ClearScreen = "\x1b[2J"
	ClearLine   = "\x1b[2K"

	HideCursor = "\x1b[?25l"
	ShowCursor = "\x1b[?25h"

	Reset = "\x1b[0m"

	Bold      = "\x1b[1m"
	Dim       = "\x1b[2m"
	Italic    = "\x1b[3m"
	Underline = "\x1b[4m"
	Blink     = "\x1b[5m"
	Reverse   = "\x1b[7m"

	NormalWeight = "\x1b[22m"
	NoItalic     = "\x1b[23m"
	NoUnderline  = "\x1b[24m"
	NoBlink      = "\x1b[25m"
	NoReverse    = "\x1b[27m"
)
