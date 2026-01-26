package core

type Style uint8

const (
	Bold Style = iota
)

type PaddingMode uint

const (
	Left PaddingMode = iota
	Right
	Center
	Fill
	Custom
)

type Padding struct {
	Left    uint16
	Right   uint16
	Padding PaddingMode
}

type Fragment struct {
	Text   string
	Styles []Style
}

type Line struct {
	Text    []Fragment
	Padding Padding
}

func NewLine(padding Padding, text string) Line {
	return Line{
		Text: []Fragment{{
			Text: text,
		}},
		Padding: padding,
	}
}

type InputLine struct {
	Prompt string
	Value  string
	Cursor int
}

type ViewModel struct {
	Lines []Line
	Input *InputLine
}
