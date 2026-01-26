package core

import "unicode/utf8"

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
	FillUp
	FillDown
	Custom
)

type Padding struct {
	Left    uint16
	Right   uint16
	Padding PaddingMode
}

func ModePadding(padding PaddingMode) Padding {
	return Padding{
		Padding: padding,
	}
}

type Fragment struct {
	Text   string
	Styles []Style
}

func NewFragment(text string, styles ...Style) Fragment {
	return Fragment{
		Text: text,
		Styles: styles,
	}
}

func UnstyledFragment(text string) Fragment {
	return NewFragment(text)
}

type Line struct {
	Text    []Fragment
	Padding Padding
}

func NewLine(text string, padding Padding) Line {
	return Line{
		Text: []Fragment{{
			Text: text,
		}},
		Padding: padding,
	}
}

func EmptyLine(padding Padding) Line {
	return Line{
		Text: []Fragment{},
		Padding: padding,
	}
}

func FragmentLine(padding Padding, fragments ...Fragment) Line {
	return Line{
		Text: fragments,
		Padding: padding,
	}
}

func (l Line) Len() int {
	lineLen := 0
	for _, v := range l.Text {
		lineLen += utf8.RuneCountInString(v.Text)
	}
	return lineLen + len(l.Text) -1 
}

type InputLine struct {
	Prompt string
	Value  string
	Cursor int
}

type ViewModel struct {
	Headers []Line
	Lines   []Line
	Input   *InputLine
}
