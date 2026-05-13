package selection

import (
	"strings"
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/model/input"
	"github.com/Rafael24595/go-reacterm-core/engine/model/offset"
	"github.com/Rafael24595/go-reacterm-core/engine/render/marker"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
)

type expectedFragment struct {
	content string
	atoms   style.Atom
}

func ef(content string, atoms ...style.Atom) expectedFragment {
	return expectedFragment{
		content: content,
		atoms:   style.MergeAtom(atoms...),
	}
}

func escapeLF(r string) string {
	return strings.ReplaceAll(r, "\n", "\\n")
}

func TestRendererResolve(t *testing.T) {
	tests := []struct {
		name string

		buffer string
		start  offset.Offset
		end    offset.Offset

		caret  offset.Offset
		anchor offset.Offset

		expectedEnd offset.Offset
		expected    []expectedFragment
	}{
		{
			name:   "forward non-enter single rune",
			buffer: "abc",
			start:  1,
			end:    2,
			caret:  2,
			anchor: 1,

			expectedEnd: 2,
			expected: []expectedFragment{
				ef("b", style.AtmFocus),
			},
		},
		{
			name:   "forward non-enter multiple runes",
			buffer: "abcd",
			start:  1,
			end:    3,
			caret:  3,
			anchor: 1,

			expectedEnd: 3,
			expected: []expectedFragment{
				ef("b"),
				ef("c", style.AtmFocus),
			},
		},
		{
			name:   "forward enter at end of buffer",
			buffer: "abc\n",
			start:  3,
			end:    4,
			caret:  4,
			anchor: 3,

			expectedEnd: 4,
			expected: []expectedFragment{
				ef(marker.PrintableCaretText),
				ef("\n"),
				ef(marker.PrintableCaretText, style.AtmFocus),
			},
		},
		{
			name:   "forward enter before visible rune",
			buffer: "abc\nd",
			start:  3,
			end:    4,
			caret:  4,
			anchor: 3,

			expectedEnd: 5,
			expected: []expectedFragment{
				ef(marker.PrintableCaretText),
				ef("\n"),
				ef("d", style.AtmFocus),
			},
		},
		{
			name:   "backward selection",
			buffer: "abcd",
			start:  1,
			end:    3,
			caret:  1,
			anchor: 3,

			expectedEnd: 3,
			expected: []expectedFragment{
				ef("bc", style.AtmFocus),
			},
		},
		{
			name:   "backward selection starting with LF",
			buffer: "a\nbc",
			start:  1,
			end:    3,
			caret:  1,
			anchor: 3,

			expectedEnd: 3,
			expected: []expectedFragment{
				ef(marker.PrintableCaretText, style.AtmFocus),
				ef("\nb"),
			},
		},
		{
			name:   "forward enter before another LF",
			buffer: "abc\n\n",
			start:  3,
			end:    4,
			caret:  4,
			anchor: 3,

			expectedEnd: 4,
			expected: []expectedFragment{
				ef(marker.PrintableCaretText),
				ef("\n"),
				ef(marker.PrintableCaretText, style.AtmFocus),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buffer := []rune(tt.buffer)

			if tt.name == "forward enter skips over visible rune after LF" {
				println()
			}

			renderer := NewRenderer(
				buffer,
				tt.start,
				tt.end,
			)

			caret := input.NewTextCursor(false)
			caret.MoveSelectTo(buffer, tt.caret, tt.anchor)

			if tt.start == tt.end {
				assert.Panic(t, func() {
					renderer.Resolve(caret)
				})
				return
			}

			result := renderer.Resolve(caret)

			assert.Equal(t, tt.expectedEnd, result.End)
			assert.Equal(t, len(tt.expected), len(result.Frags))

			for i, frag := range result.Frags {
				expected := tt.expected[i]

				assert.Equal(t, escapeLF(expected.content), escapeLF(frag.Text))
				assert.Equal(t, expected.atoms, frag.Atom)
			}
		})
	}
}
