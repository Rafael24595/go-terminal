package wrap

import (
	"strings"
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	"github.com/Rafael24595/go-reacterm-core/engine/helper/runes"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

func assembleLines(t *testing.T, lines ...text.Line) string {
	t.Helper()

	var sb strings.Builder

	for i, l := range lines {
		_, err := sb.WriteString(
			text.LineToString(&l),
		)

		assert.NotError(t, err)

		if i < len(lines)-1 {
			_, err := sb.WriteString("\n")

			assert.NotError(t, err)
		}
	}

	return sb.String()
}

func TestNormalizeLines_Integrity(t *testing.T) {
	line := text.NewLine("golang ziglang 10.50 rust")

	assert.Len(t, 1, line.Text)

	tokenized := NormalizeLines(*line)

	assert.Len(t, 1, tokenized)
	assert.Len(t, 7, tokenized[0].Text)
}

func TestWrapLine_Simple(t *testing.T) {
	line := text.NewLine(
		"HELLO WORLD",
		style.SpecFromKind(style.SpcKindPaddingLeft),
	)

	lines := Line(5, line)

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

func TestWrapLine_Styles(t *testing.T) {
	line := text.LineFromFragments(
		*text.NewFragment("HELLO").AddAtom(style.AtmBold),
		*text.NewFragment(" "),
		*text.NewFragment("WORLD"),
	).SetSpec(style.SpecFromKind(style.SpcKindPaddingLeft))

	lines := Line(7, line)

	assert.Equal(t, 2, len(lines))

	assert.Equal(t, "HELLO", lines[0].Text[0].Text)
	assert.True(t, lines[0].Text[0].Atom.HasAny(style.AtmBold))

	assert.Equal(t, " ", lines[0].Text[1].Text)

	assert.Equal(t, "WORLD", lines[1].Text[0].Text)
}

func TestWrapLine_LongWord(t *testing.T) {
	txt := "HELLO WORLD FROM GOLANG"

	line := text.NewLine(txt,
		style.SpecFromKind(style.SpcKindPaddingLeft),
	)

	maxWidth := winsize.Cols(10)
	lines := Line(maxWidth, line)

	for i, l := range lines {
		text := ""
		for _, f := range l.Text {
			text += f.Text
		}
		if runes.Measure(text) > maxWidth {
			t.Errorf("line %d too long: %s", i, text)
		}
	}

	totalRunes := winsize.Cols(0)
	for _, l := range lines {
		for _, f := range l.Text {
			totalRunes += runes.Measure(f.Text)
		}
	}
	if totalRunes != runes.Measure(txt) {
		t.Errorf("total runes mismatch")
	}
}

func TestWrapLine_MultipleFragments(t *testing.T) {
	line := text.LineFromFragments(
		*text.NewFragment("HELLO").AddAtom(style.AtmBold),
		*text.NewFragment("WORLD").AddAtom(style.AtmBold),
		*text.NewFragment("GO"),
	).SetSpec(style.SpecFromKind(style.SpcKindPaddingLeft))

	maxWidth := winsize.Cols(8)
	lines := Line(maxWidth, line)

	for _, l := range lines {
		width := winsize.Cols(0)
		for _, f := range l.Text {
			width += runes.Measure(f.Text)
		}
		if width > maxWidth {
			t.Errorf("line exceeds maxWidth: %v", l)
		}
	}
}

func TestNextWrappedLine_Fit(t *testing.T) {
	line := text.NewLine("golang")

	got, remain := NextLine(10, []text.Line{*line})

	assert.Equal(t, "golang", text.LineToString(got))

	assert.Len(t, 0, remain)
}

func TesNextWrappedLine_Split(t *testing.T) {
	line := text.NewLine("golang")

	got, remain := NextLine(2, []text.Line{*line})

	assert.Equal(t, "go", text.LineToString(got))

	assert.Len(t, 1, remain)
	assert.Equal(t, "lang", text.LineToString(&remain[0]))
}

func TesNextWrappedLine_MultiFragment(t *testing.T) {
	line := text.LineFromFragments(
		*text.NewFragment("go"),
		*text.NewFragment(" "),
		*text.NewFragment("zig"),
		*text.NewFragment(" "),
		*text.NewFragment("c++"),
	)

	got, remain := NextLine(6, []text.Line{*line})

	assert.Equal(t, "go zig", text.LineToString(got))
	assert.Len(t, 1, remain)

	assert.Equal(t, " c++", text.LineToString(&remain[0]))
}

func TesNextWrappedLine_BreakLongWordSingleFragment(t *testing.T) {
	line := text.NewLine("golangziglangrustlang")

	got, remain := NextLine(6, []text.Line{*line})
	assert.Equal(t, "golang", text.LineToString(got))

	assert.Equal(t, "ziglangrustlang", text.LineToString(&remain[0]))
}

func TesNextWrappedLine_BreakLongWordMultipleFragments(t *testing.T) {
	line := text.LineFromFragments(
		*text.NewFragment("golang"),
		*text.NewFragment(" "),
		*text.NewFragment("zigrust"),
	)

	got, remain := NextLine(10, []text.Line{*line})
	assert.Equal(t, "golang ", text.LineToString(got))

	assert.Equal(t, "zigrust", text.LineToString(&remain[0]))
}

func TestSplitLineFeeds(t *testing.T) {
	tests := []struct {
		name         string
		input        *text.Line
		expectedSize int
		expectedText string
	}{
		{
			name: "WithoutLineFeed",
			input: text.EmptyLine().PushFragments(
				*text.NewFragment("Hello Golang"),
			),
			expectedSize: 1,
			expectedText: "Hello Golang",
		},
		{
			name: "SingleLineFeed",
			input: text.EmptyLine().PushFragments(
				*text.NewFragment("Golang\nZiglang"),
			),
			expectedSize: 2,
			expectedText: "Golang\nZiglang",
		},
		{
			name: "LineFeedBetweenFragments",
			input: text.EmptyLine().PushFragments(
				*text.NewFragment("Rust"),
				*text.NewFragment("\nZig"),
			),
			expectedSize: 2,
			expectedText: "Rust\nZig",
		},
		{
			name: "MultipleLineFeedWithEmptyLine",
			input: text.EmptyLine().PushFragments(
				*text.NewFragment("Go\n\nC++"),
			),
			expectedSize: 3,
			expectedText: "Go\n\nC++",
		},
		{
			name: "LineFeedAtEnd",
			input: text.EmptyLine().PushFragments(
				*text.NewFragment("Rust\n"),
			),
			expectedSize: 2,
			expectedText: "Rust\n",
		},
		{
			name: "LineFeedWithCarriageReturn",
			input: text.EmptyLine().PushFragments(
				*text.NewFragment("Zig\r\nGolang"),
			),
			expectedSize: 2,
			expectedText: "Zig\nGolang",
		},
		{
			name: "CarriageReturn",
			input: text.EmptyLine().PushFragments(
				*text.NewFragment("Java\rElixir"),
			),
			expectedSize: 2,
			expectedText: "Java\nElixir",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := splitLineFeeds(tt.input)

			assert.Len(t, tt.expectedSize, got)
			assert.Equal(t, tt.expectedText, assembleLines(t, got...))
		})
	}
}
