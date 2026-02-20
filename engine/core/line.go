package core

import (
	"strings"
	"unicode/utf8"

	"github.com/Rafael24595/go-terminal/engine/core/style"
)

type Fragment struct {
	Text string
	Atom style.Atom
	Spec style.Spec
}

func FragmentsFromString(text ...string) []Fragment {
	fragments := make([]Fragment, len(text))
	for i, v := range text {
		fragments[i] = NewFragment(v)
	}
	return fragments
}

func NewFragment(text string) Fragment {
	return Fragment{
		Text: text,
		Atom: style.AtmNone,
		Spec: style.SpecEmpty(),
	}
}

func EmptyFragment() Fragment {
	return NewFragment("")
}

func (f Fragment) AddAtom(styles ...style.Atom) Fragment {
	f.Atom = style.MergeAtom(styles...)
	return f
}

func (f Fragment) AddSpec(styles ...style.Spec) Fragment {
	f.Spec = style.MergeSpec(styles...)
	return f
}

func (f Fragment) Len() int {
	lineLen := utf8.RuneCountInString(f.Text)
	return style.SpecLen(f.Spec, lineLen)
}

type Line struct {
	Order uint16
	Text  []Fragment
	Spec  style.Spec
}

func NewLines(lines ...Line) []Line {
	return lines
}

func FixedLinesFromLines(style style.Spec, lines ...Line) []Line {
	for i := range lines {
		lines[i].Spec = style
	}
	return lines
}

func LineFromFragments(fragments ...Fragment) Line {
	return Line{
		Text: fragments,
		Spec: style.SpecEmpty(),
	}
}

func NewLine(text string, style style.Spec) Line {
	return Line{
		Text: []Fragment{{
			Text: text,
		}},
		Spec: style,
	}
}

func LineFromString(text string, styles ...style.Atom) Line {
	return Line{
		Text: []Fragment{{
			Text: text,
			Atom: style.MergeAtom(styles...),
		}},
		Spec: style.SpecEmpty(),
	}
}

func LineFromSpec(style style.Spec) Line {
	return Line{
		Text: []Fragment{},
		Spec: style,
	}
}

func LineJump() Line {
	return Line{
		Text: FragmentsFromString(""),
		Spec: style.SpecFromKind(style.SpcKindFill),
	}
}

func FragmentLine(style style.Spec, fragments ...Fragment) Line {
	return Line{
		Text: fragments,
		Spec: style,
	}
}

func (l *Line) SetOrder(order uint16) *Line {
	l.Order = order
	return l
}

func (l Line) Len() int {
	lineLen := 0
	for _, v := range l.Text {
		lineLen += v.Len()
	}
	return lineLen
}

func (l Line) String() string {
	buffer := make([]string, 0)
	for _, v := range l.Text {
		buffer = append(buffer, v.Text)
	}
	return strings.Join(buffer, " ")
}
