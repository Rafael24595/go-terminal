package core

import (
	"strings"
	"testing"

	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/test/support/assert"
)

func tokenString(tok core.WordToken) string {
	var b strings.Builder
	for _, f := range tok.Text {
		b.WriteString(f.Text)
	}
	return b.String()
}

func tokenStrings(tokens []core.WordToken) []string {
	out := make([]string, len(tokens))
	for i, t := range tokens {
		out[i] = tokenString(t)
	}
	return out
}

func TestTokenizeLine(t *testing.T) {
	tests := []struct {
		name     string
		line     core.Line
		expected []string
	}{
		{
			name: "single word",
			line: core.LineFromFragments(
				core.FragmentsFromString("Golang")...,
			),
			expected: []string{"Golang"},
		},
		{
			name: "word split across fragments",
			line: core.LineFromFragments(
				core.FragmentsFromString("Z", "ig", "lang")...,
			),
			expected: []string{"Ziglang"},
		},
		{
			name: "two words with space",
			line: core.LineFromFragments(
				core.FragmentsFromString("hello cargo")...,
			),
			expected: []string{"hello", " ", "cargo"},
		},
		{
			name: "multiple spaces preserved",
			line: core.LineFromFragments(
				core.FragmentsFromString("hello   golangci")...,
			),
			expected: []string{"hello", "   ", "golangci"},
		},
		{
			name: "spaces across fragments",
			line: core.LineFromFragments(
				core.FragmentsFromString("hello", "  ", "zig")...,
			),
			expected: []string{"hello", "  ", "zig"},
		},
		{
			name: "styled per character",
			line: core.LineFromFragments(
				core.FragmentsFromString("r", "u", "s", "t", "c")...,
			),
			expected: []string{"rustc"},
		},
		{
			name: "leading and trailing spaces",
			line: core.LineFromFragments(
				core.FragmentsFromString("  Golang  ")...,
			),
			expected: []string{"  ", "Golang", "  "},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := core.TokenizeLineWords(tt.line)
			got := tokenStrings(tokens)

			assert.Equal(t, len(tt.expected), len(got))
			for i := range got {
				assert.Equal(t, tt.expected[i], got[i])
			}
		})
	}
}

func TestTokenizeLine_EmptyLine(t *testing.T) {
	line := core.LineFromFragments()

	tokens := core.TokenizeLineWords(line)

	assert.Equal(t, 0, len(tokens))
}

func TestTokenizeLine_EmptyFragmentIgnored(t *testing.T) {
	line := core.LineFromFragments(
		core.NewFragment(""),
		core.NewFragment("Golang"),
	)

	tokens := core.TokenizeLineWords(line)

	assert.Equal(t, 1, len(tokens))
	assert.Equal(t, "Golang", tokenString(tokens[0]))
}

func TestTokenizeLine_OnlySpaces(t *testing.T) {
	line := core.LineFromFragments(
		core.NewFragment("   "),
	)

	tokens := core.TokenizeLineWords(line)

	assert.Equal(t, 1, len(tokens))
	assert.Equal(t, "   ", tokenString(tokens[0]))
}

func TestTokenizeLine_StyleChangeRequiresFragmentSplit(t *testing.T) {
	line := core.LineFromFragments(
		core.NewFragment("Zig", core.Bold),
		core.NewFragment("lang", core.Bold),
	)

	tokens := core.TokenizeLineWords(line)

	assert.Equal(t, 1, len(tokens))
	assert.Equal(t, 2, len(tokens[0].Text))

	assert.True(t, tokens[0].Text[0].Styles.HasAny(core.Bold))
	assert.True(t, tokens[0].Text[1].Styles.HasAny(core.Bold))
}

func TestTokenizeLine_PreservesStylesAcrossFragments(t *testing.T) {
	line := core.LineFromFragments(
		core.NewFragment("ru"),
		core.NewFragment("st", core.Select),
		core.NewFragment("up"),
	)

	tokens := core.TokenizeLineWords(line)

	assert.Equal(t, 1, len(tokens))

	assert.True(t, tokens[0].Text[1].Styles.HasAny(core.Select))
}

func TestTokenizeLine_MultipleSpaceFragmentsKeepStyles(t *testing.T) {
	line := core.LineFromFragments(
		core.NewFragment(" ", core.Bold),
		core.NewFragment(" ", core.Select),
		core.NewFragment("c", core.Bold),
	)

	tokens := core.TokenizeLineWords(line)

	assert.Equal(t, 2, len(tokens))

	assert.Equal(t, 2, len(tokens[0].Text))

	assert.Equal(t, " ", tokens[0].Text[0].Text)
	assert.True(t, tokens[0].Text[0].Styles.HasAny(core.Bold))

	assert.Equal(t, " ", tokens[0].Text[1].Text)
	assert.True(t, tokens[0].Text[1].Styles.HasAny(core.Select))

	assert.Equal(t, 1, len(tokens[1].Text))
	assert.Equal(t, "c", tokens[1].Text[0].Text)
	assert.True(t, tokens[1].Text[0].Styles.HasAny(core.Bold))
}

func TestTokenizeLine_FinalFlushPreservesStyles(t *testing.T) {
	line := core.LineFromFragments(
		core.NewFragment("c++", core.Bold),
	)

	tokens := core.TokenizeLineWords(line)

	assert.Equal(t, 1, len(tokens))

	assert.True(t, tokens[0].Text[0].Styles.HasAny(core.Bold))
}

func TestSplitLongToken_PreservesStyles(t *testing.T) {
	token := core.WordTokenFromFragments(
		core.NewFragment("abcdef", core.Bold),
	)

	line, emitted, width := core.SplitLongToken(
		token, 3,
		core.LineFromPadding(core.ModePadding(core.Fill)), 0,
	)

	assert.Equal(t, 3, width)

	assert.Equal(t, 3, line.Len())
	assert.Equal(t, "def", line.String())

	assert.Equal(t, 1, len(emitted))
	assert.Equal(t, "abc", emitted[0].String())

	assert.True(t, emitted[0].Text[0].Styles.HasAny(core.Bold))
}

func TestSplitLongToken_WithInitialWidth(t *testing.T) {
	token := core.WordTokenFromFragments(
		core.NewFragment("abcdef"),
	)

	current := core.LineFromPadding(core.ModePadding(core.Fill))
	current.Text = append(current.Text, core.NewFragment("XY"))

	line, emitted, width := core.SplitLongToken(
		token, 4,
		current, 2,
	)

	assert.Equal(t, "cdef", line.String())

	assert.Equal(t, 1, len(emitted))
	assert.Equal(t, 2, len(emitted[0].Text))

	assert.Equal(t, "XY", emitted[0].Text[0].Text)
	assert.Equal(t, "ab", emitted[0].Text[1].Text)

	assert.Equal(t, 4, width)
}

func TestSplitLongToken_CurrentAlreadyFull(t *testing.T) {
	token := core.WordTokenFromFragments(
		core.NewFragment("abc"),
	)

	current := core.LineFromPadding(core.ModePadding(core.Fill))
	current.Text = append(current.Text, core.NewFragment("WXYZ"))

	line, emitted, width := core.SplitLongToken(
		token, 4,
		current, 4,
	)

	assert.Equal(t, "abc", line.String())
	assert.Equal(t, 1, len(emitted))
	assert.Equal(t, "WXYZ", emitted[0].String())
	assert.Equal(t, 3, width)
}
