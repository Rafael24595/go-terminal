package core

import (
	"unicode"
	"unicode/utf8"

	"github.com/Rafael24595/go-terminal/engine/app/state"
)
type InputLine struct {
	Prompt string
	Value  string
	Cursor int
}

type ViewModel struct {
	Header []Line
	Lines  []Line
	Footer []Line
	Input  *InputLine
	Pager  state.PagerState
	Cursor state.CursorState
}

func ViewModelFromUIState(state state.UIState) *ViewModel {
	return &ViewModel{
		Pager:  state.Pager,
		Cursor: state.Cursor,
	}
}

func (v *ViewModel) AddHeader(headers ...Line) *ViewModel {
	v.Header = append(v.Header, headers...)
	return v
}

func (v *ViewModel) AddLines(lines ...Line) *ViewModel {
	v.Lines = append(v.Lines, lines...)
	return v
}

func (v *ViewModel) AddFooter(footer []Line) *ViewModel {
	v.Footer = append(v.Footer, footer...)
	return v
}

func (v *ViewModel) SetPager(pager state.PagerState) *ViewModel {
	v.Pager = pager
	return v
}

func (v *ViewModel) SetCursor(cursor state.CursorState) *ViewModel {
	v.Cursor = cursor
	return v
}

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
		size += utf8.RuneCountInString(v.Text)
	}
	return size
}

func TokenizeLineWords(line Line) []WordToken {
	tokens := make([]WordToken, 0)
	fragments := make([]Fragment, 0)

	flush := func(pf, cf Fragment) Fragment {
		if len(cf.Text) > 0 {
			fragments = append(fragments, cf)
		}

		if len(fragments) > 0 {
			token := WordToken{
				Text: fragments,
			}

			tokens = append(tokens, token)
			fragments = make([]Fragment, 0)
		}

		return FragmentFromStyle(pf.Styles)
	}

	inSpace := false

	for _, frag := range line.Text {
		fragment := FragmentFromStyle(frag.Styles)
		for _, r := range frag.Text {
			if unicode.IsSpace(r) {
				if !inSpace {
					fragment = flush(frag, fragment)
				}

				fragment.Text += string(r)

				inSpace = true

				continue
			}

			if inSpace {
				fragment = flush(frag, fragment)
			}

			fragment.Text += string(r)

			inSpace = false
		}

		fragments = append(fragments, fragment)
	}

	if len(fragments) > 0 {
		flush(Fragment{}, Fragment{})
	}

	return tokens
}

func SplitLongToken(word WordToken, cols int, current Line, width int) (Line, []Line, int) {
	emmited := make([]Line, 0)
	fragments := append([]Fragment{}, word.Text...)

	flush := func() {
		emmited = append(emmited, current)
		current = LineFromPadding(current.Padding)
		width = 0
	}

	for len(fragments) > 0 {
		remaining := cols - width
		if remaining == 0 {
			flush()
			continue
		}

		fragment := fragments[0]
		size := utf8.RuneCountInString(fragment.Text)

		if size <= remaining {
			current.Text = append(current.Text, fragment)
			width += size

			fragments = fragments[1:]

			continue
		}

		taken, rest := takeFromFragment(fragment, remaining)

		current.Text = append(current.Text, taken)
		width += utf8.RuneCountInString(taken.Text)

		fragments[0] = rest

		flush()
	}

	return current, emmited, width
}

func takeFromFragment(f Fragment, n int) (taken Fragment, rest Fragment) {
	runes := []rune(f.Text)

	if n >= len(runes) {
		return f, FragmentFromStyle(f.Styles)
	}

	taken = FragmentFromStyle(f.Styles)
	taken.Text = string(runes[:n])

	rest = FragmentFromStyle(f.Styles)
	rest.Text = string(runes[n:])

	return
}
