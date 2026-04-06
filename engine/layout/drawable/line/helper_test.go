package line

import (
	"strings"
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-terminal/engine/helper/runes"
	"github.com/Rafael24595/go-terminal/engine/render/style"
	"github.com/Rafael24595/go-terminal/engine/render/text"
)

func TestSplitLineWords_Simple(t *testing.T) {
	line := text.NewLine(
		"HELLO WORLD",
		style.SpecFromKind(style.SpcKindPaddingLeft),
	)

	maxWidth := 5
	lines := WrapLineWords(maxWidth, line)

	expected := []string{"HELLO", " ", "WORLD"}

	assert.Len(t, len(expected), lines)

	for i, l := range lines {
		var text strings.Builder
		for _, f := range l.Text {
			text.WriteString(f.Text)
		}

		assert.Equal(t, expected[i], text.String())
	}
}

func TestSplitLineWords_Styles(t *testing.T) {
	line := text.FragmentLine(
		style.SpecFromKind(style.SpcKindPaddingLeft),
		text.NewFragment("HELLO").AddAtom(style.AtmBold),
		text.NewFragment(" "),
		text.NewFragment("WORLD"),
	)

	maxWidth := 7
	lines := WrapLineWords(maxWidth, line)

	assert.Equal(t, 2, len(lines))

	assert.Equal(t, "HELLO", lines[0].Text[0].Text)
	assert.True(t, lines[0].Text[0].Atom.HasAny(style.AtmBold))

	assert.Equal(t, " ", lines[0].Text[1].Text)

	assert.Equal(t, "WORLD", lines[1].Text[0].Text)
}

func TestSplitLineWords_LongWord(t *testing.T) {
	txt := "HELLO WORLD FROM GOLANG"

	line := text.NewLine(
		txt,
		style.SpecFromKind(style.SpcKindPaddingLeft),
	)

	maxWidth := 10
	lines := WrapLineWords(maxWidth, line)

	for i, l := range lines {
		text := ""
		for _, f := range l.Text {
			text += f.Text
		}
		if runes.Measure(text) > maxWidth {
			t.Errorf("line %d too long: %s", i, text)
		}
	}

	totalRunes := 0
	for _, l := range lines {
		for _, f := range l.Text {
			totalRunes += runes.Measure(f.Text)
		}
	}
	if totalRunes != runes.Measure(txt) {
		t.Errorf("total runes mismatch")
	}
}

func TestSplitLineWords_MultipleFragments(t *testing.T) {
	line := text.FragmentLine(
		style.SpecFromKind(style.SpcKindPaddingLeft),
		text.NewFragment("HELLO").AddAtom(style.AtmBold),
		text.NewFragment("WORLD").AddAtom(style.AtmBold),
		text.NewFragment("GO"),
	)

	maxWidth := 8
	lines := WrapLineWords(maxWidth, line)

	for _, l := range lines {
		width := 0
		for _, f := range l.Text {
			width += runes.Measure(f.Text)
		}
		if width > maxWidth {
			t.Errorf("line exceeds maxWidth: %v", l)
		}
	}
}
