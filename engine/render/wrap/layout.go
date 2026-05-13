package wrap

import "github.com/Rafael24595/go-reacterm-core/engine/render/text"

type LayoutLine struct {
	Source *text.Line
	Words  []word
}

func NewLayoutLine(source *text.Line, words ...word) *LayoutLine {
	return &LayoutLine{
		Source: source,
		Words:  words,
	}
}

func (l *LayoutLine) toLine() *text.Line {
	line := text.LineFromMeta(l.Source)
	for _, v := range l.Words {
		line.PushFragments(v.Text...)
	}
	return line
}

func (l *LayoutLine) clone() *LayoutLine {
	newLine := NewLayoutLine(l.Source)
	newLine.Words = make([]word, len(l.Words))
	copy(newLine.Words, l.Words)
	return newLine
}

func CloneLayoutLines(lines ...LayoutLine) []LayoutLine {
	clones := make([]LayoutLine, len(lines))
	for i, v := range lines {
		clones[i] = *v.clone()
	}
	return clones
}
