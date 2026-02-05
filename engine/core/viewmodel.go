package core

import (
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
