package text

import (
	"strings"

	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
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
	l.AddSpec(other.Spec)
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

func (l *Line) Clone() *Line {
	newLine := EmptyLine().CopyMeta(l)
	newLine.Text = make([]Fragment, len(l.Text))
	copy(newLine.Text, l.Text)
	return newLine
}

func LineMeasure(line *Line, cols int) int {
	return style.SpecMeasure(line.Spec, style.LayoutContext{
		Text: FragmentMeasure(cols, line.Text...),
		Cols: cols,
	})
}

func LineToString(line *Line) string {
	buffer := make([]string, 0)
	for _, v := range line.Text {
		buffer = append(buffer, v.Text)
	}
	return strings.Join(buffer, "")
}
