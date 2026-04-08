package text

import (
	"strings"
	"unicode"

	"github.com/Rafael24595/go-terminal/engine/helper/runes"
	"github.com/Rafael24595/go-terminal/engine/render/style"
)

type WordToken struct {
	Text []Fragment
}

func WordTokenFromFragments(fragments ...Fragment) WordToken {
	return WordToken{
		Text: fragments,
	}
}

func (t WordToken) Size() int {
	size := 0
	for _, v := range t.Text {
		size += FragmentMeasure(v)
	}
	return size
}

func TokenizeLineWords(line Line) []WordToken {
	tokens := make([]WordToken, 0, len(line.Text))
	fragments := make([]Fragment, 0, 4)

	var sb strings.Builder

	flush := func(frag Fragment) {
		if sb.Len() > 0 {
			f := EmptyFragmentFrom(frag)
			f.Text = sb.String()
			fragments = append(fragments, f)
			sb.Reset()
		}

		if len(fragments) > 0 {
			tokenFrags := make([]Fragment, len(fragments))
			copy(tokenFrags, fragments)

			token := WordToken{
				Text: tokenFrags,
			}

			tokens = append(tokens, token)
			fragments = fragments[:0]
		}
	}

	inSpace := false

	for _, frag := range line.Text {
		if frag.Atom.HasAny(style.AtmWrap) {
			tokens = append(tokens, WordToken{
				Text: []Fragment{frag},
			})

			continue
		}

		for _, r := range frag.Text {
			isSpace := unicode.IsSpace(r)

			if isSpace != inSpace {
				flush(frag)
			}

			inSpace = isSpace
			sb.WriteRune(r)
		}

		if sb.Len() > 0 {
			f := EmptyFragmentFrom(frag)
			f.Text = sb.String()
			fragments = append(fragments, f)
			sb.Reset()
		}
	}

	if len(fragments) > 0 {
		flush(Fragment{})
	}

	return tokens
}

func SplitLongToken(word WordToken, cols int, current Line, width int) (Line, []Line, int) {
	emmited := make([]Line, 0)
	if cols <= 0 {
		emmited = append(emmited, LineFromFragments(word.Text...))
		return current, emmited, 0
	}

	fragments := word.Text

	flush := func() {
		emmited = append(emmited, current)
		current = LineFromSpec(current.Spec)
		width = 0
	}

	for len(fragments) > 0 {
		remaining := cols - width
		if remaining == 0 {
			flush()
			continue
		}

		fragment := fragments[0]
		size := FragmentMeasure(fragment)

		if size <= remaining {
			current.Text = append(current.Text, fragment)
			width += size

			fragments = fragments[1:]

			continue
		}

		taken, rest := takeFromFragment(fragment, remaining)

		current.Text = append(current.Text, taken)
		width += FragmentMeasure(taken)

		fragments = append([]Fragment{rest}, fragments[1:]...)

		flush()
	}

	return current, emmited, width
}

func takeFromFragment(f Fragment, n int) (Fragment, Fragment) {
	if n <= 0 {
		return EmptyFragmentFrom(f), f
	}

	byteIndex, canBreak := runes.RuneIndexToByteIndex(f.Text, n)
	if !canBreak {
		return f, EmptyFragmentFrom(f)
	}

	taken := EmptyFragmentFrom(f)
	taken.Text = f.Text[:byteIndex]

	rest := EmptyFragmentFrom(f)
	rest.Text = f.Text[byteIndex:]

	return taken, rest
}
