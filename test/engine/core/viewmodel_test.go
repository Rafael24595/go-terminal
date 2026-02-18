package core_test

import (
	"strings"
	"testing"

	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/core/style"
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
		core.NewFragment("Zig").AddAtom(style.AtmBold),
		core.NewFragment("lang").AddAtom(style.AtmBold),
	)

	tokens := core.TokenizeLineWords(line)

	assert.Equal(t, 1, len(tokens))
	assert.Equal(t, 2, len(tokens[0].Text))

	assert.True(t, tokens[0].Text[0].Atom.HasAny(style.AtmBold))
	assert.True(t, tokens[0].Text[1].Atom.HasAny(style.AtmBold))
}

func TestTokenizeLine_PreservesStylesAcrossFragments(t *testing.T) {
	line := core.LineFromFragments(
		core.NewFragment("ru"),
		core.NewFragment("st").AddAtom(style.AtmSelect),
		core.NewFragment("up"),
	)

	tokens := core.TokenizeLineWords(line)

	assert.Equal(t, 1, len(tokens))

	assert.True(t, tokens[0].Text[1].Atom.HasAny(style.AtmSelect))
}

func TestTokenizeLine_MultipleSpaceFragmentsKeepStyles(t *testing.T) {
	line := core.LineFromFragments(
		core.NewFragment(" ").AddAtom(style.AtmBold),
		core.NewFragment(" ").AddAtom(style.AtmSelect),
		core.NewFragment("c").AddAtom(style.AtmBold),
	)

	tokens := core.TokenizeLineWords(line)

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
	line := core.LineFromFragments(
		core.NewFragment("c++").AddAtom(style.AtmBold),
	)

	tokens := core.TokenizeLineWords(line)

	assert.Equal(t, 1, len(tokens))

	assert.True(t, tokens[0].Text[0].Atom.HasAny(style.AtmBold))
}

func TestSplitLongToken_PreservesStyles(t *testing.T) {
	token := core.WordTokenFromFragments(
		core.NewFragment("abcdef").AddAtom(style.AtmBold),
	)

	line, emitted, width := core.SplitLongToken(
		token, 3,
		core.LineFromSpec(style.SpecFromKind(style.SpcKindFill)), 0,
	)

	assert.Equal(t, 3, width)

	assert.Equal(t, 3, line.Len())
	assert.Equal(t, "def", line.String())

	assert.Equal(t, 1, len(emitted))
	assert.Equal(t, "abc", emitted[0].String())

	assert.True(t, emitted[0].Text[0].Atom.HasAny(style.AtmBold))
}

func TestSplitLongToken_WithInitialWidth(t *testing.T) {
	token := core.WordTokenFromFragments(
		core.NewFragment("abcdef"),
	)

	current := core.LineFromSpec(style.SpecFromKind(style.SpcKindFill))
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

	current := core.LineFromSpec(style.SpecFromKind(style.SpcKindFill))
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
