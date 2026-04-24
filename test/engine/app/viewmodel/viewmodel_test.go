package viewmodel_test

import (
	"strings"
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

func tokenString(token text.WordToken) string {
	var b strings.Builder
	for _, f := range token.Text {
		b.WriteString(f.Text)
	}
	return b.String()
}

func tokenStrings(tokens []text.WordToken) []string {
	out := make([]string, len(tokens))
	for i, t := range tokens {
		out[i] = tokenString(t)
	}
	return out
}

func TestTokenizeLine(t *testing.T) {
	tests := []struct {
		name     string
		line     *text.Line
		expected []string
	}{
		{
			name: "single word",
			line: text.LineFromFragments(
				text.FragmentsFromString("Golang")...,
			),
			expected: []string{"Golang"},
		},
		{
			name: "word split across fragments",
			line: text.LineFromFragments(
				text.FragmentsFromString("Z", "ig", "lang")...,
			),
			expected: []string{"Ziglang"},
		},
		{
			name: "two words with space",
			line: text.LineFromFragments(
				text.FragmentsFromString("hello cargo")...,
			),
			expected: []string{"hello", " ", "cargo"},
		},
		{
			name: "multiple spaces preserved",
			line: text.LineFromFragments(
				text.FragmentsFromString("hello   golangci")...,
			),
			expected: []string{"hello", "   ", "golangci"},
		},
		{
			name: "spaces across fragments",
			line: text.LineFromFragments(
				text.FragmentsFromString("hello", "  ", "zig")...,
			),
			expected: []string{"hello", "  ", "zig"},
		},
		{
			name: "styled per character",
			line: text.LineFromFragments(
				text.FragmentsFromString("r", "u", "s", "t", "c")...,
			),
			expected: []string{"rustc"},
		},
		{
			name: "leading and trailing spaces",
			line: text.LineFromFragments(
				text.FragmentsFromString("  Golang  ")...,
			),
			expected: []string{"  ", "Golang", "  "},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := text.TokenizeLineWords(tt.line)
			got := tokenStrings(tokens)

			assert.Len(t, len(tt.expected), got)
			for i := range got {
				assert.Equal(t, tt.expected[i], got[i])
			}
		})
	}
}

func TestTokenizeLine_EmptyLine(t *testing.T) {
	line := text.LineFromFragments()

	tokens := text.TokenizeLineWords(line)

	assert.Equal(t, 0, len(tokens))
}

func TestTokenizeLine_EmptyFragmentIgnored(t *testing.T) {
	line := text.LineFromFragments(
		*text.NewFragment(""),
		*text.NewFragment("Golang"),
	)

	tokens := text.TokenizeLineWords(line)

	assert.Equal(t, 1, len(tokens))
	assert.Equal(t, "Golang", tokenString(tokens[0]))
}

func TestTokenizeLine_OnlySpaces(t *testing.T) {
	line := text.LineFromFragments(
		*text.NewFragment("   "),
	)

	tokens := text.TokenizeLineWords(line)

	assert.Equal(t, 1, len(tokens))
	assert.Equal(t, "   ", tokenString(tokens[0]))
}

func TestTokenizeLine_StyleChangeRequiresFragmentSplit(t *testing.T) {
	line := text.LineFromFragments(
		*text.NewFragment("Zig").AddAtom(style.AtmBold),
		*text.NewFragment("lang").AddAtom(style.AtmBold),
	)

	tokens := text.TokenizeLineWords(line)

	assert.Equal(t, 1, len(tokens))
	assert.Equal(t, 2, len(tokens[0].Text))

	assert.True(t, tokens[0].Text[0].Atom.HasAny(style.AtmBold))
	assert.True(t, tokens[0].Text[1].Atom.HasAny(style.AtmBold))
}

func TestTokenizeLine_PreservesStylesAcrossFragments(t *testing.T) {
	line := text.LineFromFragments(
		*text.NewFragment("ru"),
		*text.NewFragment("st").AddAtom(style.AtmSelect),
		*text.NewFragment("up"),
	)

	tokens := text.TokenizeLineWords(line)

	assert.Equal(t, 1, len(tokens))

	assert.True(t, tokens[0].Text[1].Atom.HasAny(style.AtmSelect))
}

func TestTokenizeLine_MultipleSpaceFragmentsKeepStyles(t *testing.T) {
	line := text.LineFromFragments(
		*text.NewFragment(" ").AddAtom(style.AtmBold),
		*text.NewFragment(" ").AddAtom(style.AtmSelect),
		*text.NewFragment("c").AddAtom(style.AtmBold),
	)

	tokens := text.TokenizeLineWords(line)

	assert.Equal(t, 2, len(tokens))

	assert.Equal(t, 2, len(tokens[0].Text))

	assert.Equal(t, " ", tokens[0].Text[0].Text)
	assert.True(t, tokens[0].Text[0].Atom.HasAny(style.AtmBold))

	assert.Equal(t, " ", tokens[0].Text[1].Text)
	assert.True(t, tokens[0].Text[1].Atom.HasAny(style.AtmSelect))

	assert.Equal(t, 1, len(tokens[1].Text))
	assert.Equal(t, "c", tokens[1].Text[0].Text)
	assert.True(t, tokens[1].Text[0].Atom.HasAny(style.AtmBold))
}

func TestTokenizeLine_FinalFlushPreservesStyles(t *testing.T) {
	line := text.LineFromFragments(
		*text.NewFragment("c++").AddAtom(style.AtmBold),
	)

	tokens := text.TokenizeLineWords(line)

	assert.Equal(t, 1, len(tokens))

	assert.True(t, tokens[0].Text[0].Atom.HasAny(style.AtmBold))
}

func TestSplitLongToken_PreservesStyles(t *testing.T) {
	token := text.WordTokenFromFragments(
		*text.NewFragment("abcdef").AddAtom(style.AtmBold),
	)

	line, emitted, width := text.SplitLongToken(
		token, 3,
		*text.EmptyLine().AddSpec(style.SpecFromKind(style.SpcKindFill)), 0,
	)

	assert.Equal(t, 3, width)

	assert.Equal(t, 3, text.FragmentMeasure(3, line.Text...))
	assert.Equal(t, "def", text.LineToString(&line))

	assert.Equal(t, 1, len(emitted))
	assert.Equal(t, "abc", text.LineToString(&emitted[0]))

	assert.True(t, emitted[0].Text[0].Atom.HasAny(style.AtmBold))
}

func TestSplitLongToken_WithInitialWidth(t *testing.T) {
	token := text.WordTokenFromFragments(
		*text.NewFragment("abcdef"),
	)

	current := text.EmptyLine().AddSpec(style.SpecFromKind(style.SpcKindFill))
	current.Text = append(current.Text, *text.NewFragment("XY"))

	line, emitted, width := text.SplitLongToken(
		token, 4,
		*current, 2,
	)

	assert.Equal(t, "cdef", text.LineToString(&line))

	assert.Equal(t, 1, len(emitted))
	assert.Equal(t, 2, len(emitted[0].Text))

	assert.Equal(t, "XY", emitted[0].Text[0].Text)
	assert.Equal(t, "ab", emitted[0].Text[1].Text)

	assert.Equal(t, 4, width)
}

func TestSplitLongToken_CurrentAlreadyFull(t *testing.T) {
	token := text.WordTokenFromFragments(
		*text.NewFragment("abc"),
	)

	current := text.EmptyLine().AddSpec(style.SpecFromKind(style.SpcKindFill))
	current.Text = append(current.Text, *text.NewFragment("WXYZ"))

	line, emitted, width := text.SplitLongToken(
		token, 4,
		*current, 4,
	)

	assert.Equal(t, "abc", text.LineToString(&line))
	assert.Equal(t, 1, len(emitted))
	assert.Equal(t, "WXYZ", text.LineToString(&emitted[0]))
	assert.Equal(t, 3, width)
}
