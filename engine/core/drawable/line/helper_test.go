package line

import (
	"testing"
	"unicode/utf8"

	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/test/support/assert"
)

func TestSplitLineWords_Simple(t *testing.T) {
	line := core.NewLine(
		"HELLO WORLD",
		core.ModePadding(core.Left),
	)

	maxWidth := 5
	lines := WrapLineWords(maxWidth, line)

	expected := []string{"HELLO", " ", "WORLD"}

	assert.Equal(t, len(expected), len(lines))

	for i, l := range lines {
		text := ""
		for _, f := range l.Text {
			text += f.Text
		}

		assert.Equal(t, expected[i], text)
	}
}

func TestSplitLineWords_Styles(t *testing.T) {
	line := core.FragmentLine(
		core.ModePadding(core.Left),
		core.NewFragment("HELLO", core.Bold),
		core.NewFragment(" "),
		core.NewFragment("WORLD"),
	)

	maxWidth := 7
	lines := WrapLineWords(maxWidth, line)

	assert.Equal(t, 2, len(lines))

	assert.Equal(t, "HELLO", lines[0].Text[0].Text)
	assert.True(t, lines[0].Text[0].Styles.HasAny(core.Bold))

	assert.Equal(t, " ", lines[0].Text[1].Text)

	assert.Equal(t, "WORLD", lines[1].Text[0].Text)
}

func TestSplitLineWords_LongWord(t *testing.T) {
	text := "HELLO WORLD FROM GOLANG"

	line := core.NewLine(
		text,
		core.ModePadding(core.Left),
	)

	maxWidth := 10
	lines := WrapLineWords(maxWidth, line)

	for i, l := range lines {
		text := ""
		for _, f := range l.Text {
			text += f.Text
		}
		if utf8.RuneCountInString(text) > maxWidth {
			t.Errorf("line %d too long: %s", i, text)
		}
	}

	totalRunes := 0
	for _, l := range lines {
		for _, f := range l.Text {
			totalRunes += utf8.RuneCountInString(f.Text)
		}
	}
	if totalRunes != utf8.RuneCountInString(text) {
		t.Errorf("total runes mismatch")
	}
}

func TestSplitLineWords_MultipleFragments(t *testing.T) {
	line := core.FragmentLine(
		core.ModePadding(core.Left),
		core.NewFragment("HELLO", core.Bold),
		core.NewFragment("WORLD", core.Bold),
		core.NewFragment("GO"),
	)

	maxWidth := 8
	lines := WrapLineWords(maxWidth, line)

	for _, l := range lines {
		width := 0
		for _, f := range l.Text {
			width += utf8.RuneCountInString(f.Text)
		}
		if width > maxWidth {
			t.Errorf("line exceeds maxWidth: %v", l)
		}
	}
}
