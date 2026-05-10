package wrap

import (
	"strings"
	"unicode"

	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

type word struct {
	Text []text.Fragment
}

func splitLineWords(line *text.Line) []word {
	tokens := make([]word, 0, len(line.Text))
	frags := make([]text.Fragment, 0, 4)

	var sb strings.Builder

	flush := func(frag text.Fragment) {
		if sb.Len() > 0 {
			f := text.NewFragment(sb.String()).
				CopyMeta(&frag)
			frags = append(frags, *f)
			sb.Reset()
		}

		if len(frags) > 0 {
			tokenFrags := make([]text.Fragment, len(frags))
			copy(tokenFrags, frags)

			token := word{
				Text: tokenFrags,
			}

			tokens = append(tokens, token)
			frags = frags[:0]
		}
	}

	inSpace := false

	for _, frag := range line.Text {
		if frag.Atom.HasAny(style.AtmWrap) || text.IsStructuralFragment(frag) {
			tokens = append(tokens, word{
				Text: []text.Fragment{frag},
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
			f := text.NewFragment(sb.String()).
				CopyMeta(&frag)
			frags = append(frags, *f)
			sb.Reset()
		}
	}

	if len(frags) > 0 {
		flush(text.Fragment{})
	}

	return tokens
}

func splitLongWord(
	word word,
	cols winsize.Cols,
	current text.Line,
	width winsize.Cols,
) (text.Line, []text.Line, winsize.Cols) {
	emmited := make([]text.Line, 0)
	if cols <= 0 {
		emmited = append(emmited, *text.LineFromFragments(word.Text...))
		return current, emmited, 0
	}

	frags := word.Text

	flush := func() {
		emmited = append(emmited, current)
		current = *text.EmptyLine().AddSpec(current.Spec)
		width = 0
	}

	for len(frags) > 0 {
		remaining := cols.Clamp(width)
		if remaining == 0 {
			flush()
			continue
		}

		frag := frags[0]
		size := text.FragmentMeasure(cols, frag)

		if size <= remaining {
			current.Text = append(current.Text, frag)
			width += size

			frags = frags[1:]

			continue
		}

		taken, rest := splitFragmentAt(&frag, remaining)

		current.Text = append(current.Text, *taken)
		width += text.FragmentMeasure(cols, *taken)

		frags = append([]text.Fragment{*rest}, frags[1:]...)

		flush()
	}

	return current, emmited, width
}
