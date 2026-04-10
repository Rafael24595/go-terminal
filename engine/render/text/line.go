package text

import (
	"strings"

	"github.com/Rafael24595/go-terminal/engine/render/style"
)

type Line struct {
	Order uint16
	Text  []Fragment
	Spec  style.Spec
}

func NewLine(text string, styles ...style.Spec) *Line {
	return &Line{
		Text: []Fragment{{
			Text: text,
		}},
		Spec: style.MergeSpec(styles...),
	}
}

func EmptyLine() *Line {
	return LineFromFragments()
}

func LineFromFragments(frags ...Fragment) *Line {
	return &Line{
		Text: frags,
		Spec: style.SpecEmpty(),
	}
}

func (l *Line) CopyMeta(other *Line) *Line {
	l.Order = other.Order
	l.Spec = other.Spec
	return l
}

func (l *Line) SetOrder(order uint16) *Line {
	l.Order = order
	return l
}

func (l *Line) UnshiftFragments(frags ...Fragment) *Line {
	l.Text = append(frags, l.Text...)
	return l
}

func (l *Line) PushFragments(frags ...Fragment) *Line {
	l.Text = append(l.Text, frags...)
	return l
}

func (l *Line) AddSpec(styles ...style.Spec) *Line {
	newSpec := style.MergeSpec(styles...)
	l.Spec = style.MergeSpec(l.Spec, newSpec)
	return l
}

func (l *Line) SetSpec(styles ...style.Spec) *Line {
	l.Spec = style.MergeSpec(styles...)
	return l
}

func (l *Line) CutSpec(styles style.SpecKind) *Line {
	l.Spec, _ = style.EraseSpec(l.Spec, styles)
	return l
}

func LineFragmentsMeasure(line *Line) int {
	fragsLen := 0
	for _, f := range line.Text {
		fragsLen += FragmentMeasure(&f)
	}
	return fragsLen
}

func LineFragmentsMeasurWithContext(line *Line, ctx style.LayoutContext) int {
	fragsLen := 0
	for _, f := range line.Text {
		fragsLen += FragmentMeasureWithContext(&f, ctx)
	}
	return fragsLen
}

func LineMeasure(line *Line) int {
	fragsLen := LineFragmentsMeasure(line)
	return style.SpecMeasure(line.Spec, fragsLen)
}

func LineMeasureWithContext(line *Line, ctx style.LayoutContext) int {
	fragsLen := LineFragmentsMeasure(line)
	return style.SpecMeasureWithContext(line.Spec, fragsLen, ctx)
}

func LineToString(line *Line) string {
	buffer := make([]string, 0)
	for _, v := range line.Text {
		buffer = append(buffer, v.Text)
	}
	return strings.Join(buffer, "")
}
