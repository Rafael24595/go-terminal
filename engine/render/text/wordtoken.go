package text

import (
	"strings"
	"unicode"

	"github.com/Rafael24595/go-reacterm-core/engine/helper/runes"
	"github.com/Rafael24595/go-reacterm-core/engine/model/offset"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
)

type WordToken struct {
	Text []Fragment
}

func WordTokenFromFragments(fragments ...Fragment) WordToken {
	return WordToken{
		Text: fragments,
	}
}

func TokenizeLineWords(line *Line) []WordToken {
	tokens := make([]WordToken, 0, len(line.Text))
	frags := make([]Fragment, 0, 4)

	var sb strings.Builder

	flush := func(frag Fragment) {
		if sb.Len() > 0 {
			f := NewFragment(sb.String()).
				CopyMeta(&frag)
			frags = append(frags, *f)
			sb.Reset()
		}

		if len(frags) > 0 {
			tokenFrags := make([]Fragment, len(frags))
			copy(tokenFrags, frags)

			token := WordToken{
				Text: tokenFrags,
			}

			tokens = append(tokens, token)
			frags = frags[:0]
		}
	}

	inSpace := false

	for _, frag := range line.Text {
		if frag.Atom.HasAny(style.AtmWrap) || IsStructuralFragment(frag) {
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
			f := NewFragment(sb.String()).
				CopyMeta(&frag)
			frags = append(frags, *f)
			sb.Reset()
		}
	}

	if len(frags) > 0 {
		flush(Fragment{})
	}

	return tokens
}

func SplitLongToken(
	word WordToken,
	cols winsize.Cols,
	current Line,
	width winsize.Cols,
) (Line, []Line, winsize.Cols) {
	emmited := make([]Line, 0)
	if cols <= 0 {
		emmited = append(emmited, *LineFromFragments(word.Text...))
		return current, emmited, 0
	}

	frags := word.Text

	flush := func() {
		emmited = append(emmited, current)
		current = *EmptyLine().AddSpec(current.Spec)
		width = 0
	}

	for len(frags) > 0 {
		remaining := cols.Clamp(width)
		if remaining == 0 {
			flush()
			continue
		}

		frag := frags[0]
		size := FragmentMeasure(cols, frag)

		if size <= remaining {
			current.Text = append(current.Text, frag)
			width += size

			frags = frags[1:]

			continue
		}

		taken, rest := TakeFromFragment(&frag, remaining)

		current.Text = append(current.Text, *taken)
		width += FragmentMeasure(cols, *taken)

		frags = append([]Fragment{*rest}, frags[1:]...)

		flush()
	}

	return current, emmited, width
}

func TakeFromFragment(frag *Fragment, cols winsize.Cols) (*Fragment, *Fragment) {
	if cols <= 0 {
		return EmptyFragment().
			CopyMeta(frag), frag
	}

	byteIndex, canBreak := runes.RuneIndexToByteIndex(frag.Text, offset.Offset(cols))
	if !canBreak {
		return frag, EmptyFragment().
			CopyMeta(frag)
	}

	taken := NewFragment(frag.Text[:byteIndex]).
		CopyMeta(frag)

	rest := NewFragment(frag.Text[byteIndex:]).
		CopyMeta(frag)

	return taken, rest
}
