package core

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/Rafael24595/go-terminal/engine/app/state"
)

type Style uint8

const (
	Bold Style = iota
	Select
)

type PaddingMode uint

const (
	Left PaddingMode = iota
	Right
	Center
	Fill
	FillUp
	FillDown
	Custom
	Unstyled
)

type Padding struct {
	Left    uint16
	Right   uint16
	Padding PaddingMode
}

func ModePadding(padding PaddingMode) Padding {
	return Padding{
		Padding: padding,
	}
}

func CustomPadding(left, right uint16) Padding {
	return Padding{
		Left:    left,
		Right:   right,
		Padding: Custom,
	}
}

type Fragment struct {
	Text   string
	Styles []Style
}

func FragmentsFromString(text ...string) []Fragment {
	fragments := make([]Fragment, len(text))
	for i, v := range text {
		fragments[i] = NewFragment(v)
	}
	return fragments
}

func FragmentFromStyle(styles ...Style) Fragment {
	return Fragment{
		Text:   "",
		Styles: styles,
	}
}

func NewFragment(text string, styles ...Style) Fragment {
	return Fragment{
		Text:   text,
		Styles: styles,
	}
}

type Line struct {
	Text    []Fragment
	Padding Padding
}

func NewLines(lines ...Line) []Line {
	return lines
}

func FixedLinesFromLines(padding Padding, lines ...Line) []Line {
	for i := range lines {
		lines[i].Padding = padding
	}
	return lines
}

func LineFromFragments(fragments ...Fragment) Line {
	return Line{
		Text:    fragments,
		Padding: ModePadding(Unstyled),
	}
}

func NewLine(text string, padding Padding) Line {
	return Line{
		Text: []Fragment{{
			Text: text,
		}},
		Padding: padding,
	}
}

func LineFromString(text string) Line {
	return Line{
		Text: []Fragment{{
			Text: text,
		}},
		Padding: ModePadding(Unstyled),
	}
}

func LineFromPadding(padding Padding) Line {
	return Line{
		Text:    []Fragment{},
		Padding: padding,
	}
}

func LineJump() Line {
	return Line{
		Text:    FragmentsFromString(""),
		Padding: ModePadding(Fill),
	}
}

func FragmentLine(padding Padding, fragments ...Fragment) Line {
	return Line{
		Text:    fragments,
		Padding: padding,
	}
}

func (l Line) Len() int {
	lineLen := 0
	for _, v := range l.Text {
		lineLen += utf8.RuneCountInString(v.Text)
	}
	return lineLen
}

func (l Line) String() string {
	buffer := make([]string, 0)
	for _, v := range l.Text {
		buffer = append(buffer, v.Text)
	}
	return strings.Join(buffer, " ")
}

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

		return FragmentFromStyle(pf.Styles...)
	}

	inSpace := false

	for _, frag := range line.Text {
		fragment := FragmentFromStyle(frag.Styles...)
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
		return f, FragmentFromStyle(f.Styles...)
	}

	taken = FragmentFromStyle(f.Styles...)
	taken.Text = string(runes[:n])

	rest = FragmentFromStyle(f.Styles...)
	rest.Text = string(runes[n:])

	return
}
