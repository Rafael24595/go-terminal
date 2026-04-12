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
	line := text.LineFromFragments(
		*text.NewFragment("HELLO").AddAtom(style.AtmBold),
		*text.NewFragment(" "),
		*text.NewFragment("WORLD"),
	).SetSpec(style.SpecFromKind(style.SpcKindPaddingLeft))

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
	line := text.LineFromFragments(
		*text.NewFragment("HELLO").AddAtom(style.AtmBold),
		*text.NewFragment("WORLD").AddAtom(style.AtmBold),
		*text.NewFragment("GO"),
	).SetSpec(style.SpecFromKind(style.SpcKindPaddingLeft))

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

func TestTokenizeLines_Integrity(t *testing.T) {
    line := text.NewLine("golang ziglang 10.50 rust")
    
    assert.Len(t, 1, line.Text)

    tokenized := TokenizeLines(*line)

	assert.Len(t, 1, tokenized)
	assert.Len(t, 7, tokenized[0].Text)
}

func TestWrapNextLine_Fit(t *testing.T) {
	line := text.NewLine("golang")

	got, remain, hasMore := WrapNextLine(10, []text.Line{*line}, nil)

	assert.Equal(t, "golang", text.LineToString(got))
	assert.False(t, hasMore)

	assert.Len(t, 0, remain)
}

func TestWrapNextLine_FitWithMeta(t *testing.T) {
	meta := &IndexMeta{
		totalWidth: 4,
		sufix:      " | ",
		digits:     1,
		prefixBody: "    ",
	}

	line := text.NewLine("golang").SetOrder(1)

	got, remain, hasMore := WrapNextLine(10, []text.Line{*line}, meta)

	assert.Equal(t, "1 | golang", text.LineToString(got))
	assert.False(t, hasMore)

	assert.Len(t, 0, remain)
}

func TestWrapNextLine_Split(t *testing.T) {
	line := text.NewLine("golang")

	got, remain, hasMore := WrapNextLine(2, []text.Line{*line}, nil)

	assert.Equal(t, "go", text.LineToString(got))
	assert.True(t, hasMore)

	assert.Len(t, 1, remain)
	assert.Equal(t, "lang", text.LineToString(&remain[0]))
}

func TestWrapNextLine_SplitWithMeta(t *testing.T) {
	meta := &IndexMeta{
		totalWidth: 4,
		sufix:      " | ",
		digits:     1,
		prefixBody: " ",
	}

	line := text.NewLine("golang rust").SetOrder(1)

	got, remain, hasMore := WrapNextLine(10, []text.Line{*line}, meta)

	assert.Equal(t, "1 | golang", text.LineToString(got))
	assert.True(t, hasMore)
	assert.Len(t, 1, remain)

	got, remain, hasMore = WrapNextLine(10, remain, meta)

	assert.Equal(t, "  |  rust", text.LineToString(got))
	assert.False(t, hasMore)
	assert.Len(t, 0, remain)
}

func TestWrapNextLine_MultiFragment(t *testing.T) {
	line := text.LineFromFragments(
		*text.NewFragment("go"),
		*text.NewFragment(" "),
		*text.NewFragment("zig"),
		*text.NewFragment(" "),
		*text.NewFragment("c++"),
	)

	got, remain, hasMore := WrapNextLine(6, []text.Line{*line}, nil)

	assert.Equal(t, "go zig", text.LineToString(got))
	assert.True(t, hasMore)
	assert.Len(t, 1, remain)

	assert.Equal(t, " c++", text.LineToString(&remain[0]))
}

func TestWrapNextLine_BreakLongWordSingleFragment(t *testing.T) {
	line := text.NewLine("golangziglangrustlang")

	got, remain, hasMore := WrapNextLine(6, []text.Line{*line}, nil)
	assert.Equal(t, "golang", text.LineToString(got))
	assert.True(t, hasMore)

	assert.Equal(t, "ziglangrustlang", text.LineToString(&remain[0]))
}

func TestWrapNextLine_BreakLongWordMultipleFragments(t *testing.T) {
	line := text.LineFromFragments(
		*text.NewFragment("golang"),
		*text.NewFragment(" "),
		*text.NewFragment("zigrust"),
	)

	got, remain, hasMore := WrapNextLine(10, []text.Line{*line}, nil)
	assert.Equal(t, "golang zig", text.LineToString(got))
	assert.True(t, hasMore)

	assert.Equal(t, "rust", text.LineToString(&remain[0]))
}

func TestWrapNextLine_IndexShouldBeLesser(t *testing.T) {
	meta := &IndexMeta{
		totalWidth: 4,
		sufix:      " | ",
		digits:     1,
		prefixBody: " ",
	}

	line := text.NewLine("golang").SetOrder(1)

	assert.Panic(t, func() {
		WrapNextLine(4, []text.Line{*line}, meta)
	})
}
