package core

import (
	"strings"
	"unicode/utf8"
)

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
	Unstyled
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

func CustomPadding(left, right uint16) Padding {
	return Padding{
		Left:    left,
		Right:   right,
		Padding: Custom,
	}
}

type Fragment struct {
	Text   string
	Styles []Style
}

func FragmentsFromString(text ...string) []Fragment {
	fragments := make([]Fragment, len(text))
	for i, v := range text {
		fragments[i] = NewFragment(v)
	}
	return fragments
}

func NewFragment(text string, styles ...Style) Fragment {
	return Fragment{
		Text:   text,
		Styles: styles,
	}
}

type Line struct {
	Text    []Fragment
	Padding Padding
}

func NewLines(lines ...Line) []Line {
	return lines
}

func FixedLinesFromLines(padding Padding, lines ...Line) []Line {
	for i := range lines {
		lines[i].Padding = padding
	}
	return lines
}

func NewLine(text string, padding Padding) Line {
	return Line{
		Text: []Fragment{{
			Text: text,
		}},
		Padding: padding,
	}
}

func LineFromString(text string) Line {
	return Line{
		Text: []Fragment{{
			Text: text,
		}},
		Padding: ModePadding(Unstyled),
	}
}

func LineFromPadding(padding Padding) Line {
	return Line{
		Text:    []Fragment{},
		Padding: padding,
	}
}

func LineJump() Line {
	return Line{
		Text:    FragmentsFromString(""),
		Padding: ModePadding(Fill),
	}
}

func FragmentLine(padding Padding, fragments ...Fragment) Line {
	return Line{
		Text:    fragments,
		Padding: padding,
	}
}

func (l Line) Len() int {
	lineLen := 0
	for _, v := range l.Text {
		lineLen += utf8.RuneCountInString(v.Text)
	}
	return lineLen + len(l.Text) - 1
}

func (l Line) String() string {
	buffer := make([]string, 0)
	for _, v := range l.Text {
		buffer = append(buffer, v.Text)
	}
	return strings.Join(buffer, " ")
}

type InputLine struct {
	Prompt string
	Value  string
	Cursor int
}

type ViewModel struct {
	Cursor *uint
	Page   *uint
	Header []Line
	Lines  []Line
	Footer []Line
	Input  *InputLine
}
