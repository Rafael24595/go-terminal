package wrapper_layout

import (
	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/app/viewmodel"
	"github.com/Rafael24595/go-terminal/engine/layout"
	"github.com/Rafael24595/go-terminal/engine/render/text"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

type FixedLayout struct {
	layout  layout.Layout
	maxRows uint16
	maxCols uint16
}

func NewFixed(layout layout.Layout, maxRows, maxCols uint16) *FixedLayout {
	return &FixedLayout{
		layout:  layout,
		maxRows: maxRows,
		maxCols: maxCols,
	}
}

func (l *FixedLayout) Update(maxRows uint16, maxCols uint16) *FixedLayout {
	l.maxCols = maxCols
	l.maxRows = maxRows
	return l
}

func (l *FixedLayout) ToLayout() layout.Layout {
	return layout.Layout{
		Apply: l.Appy,
	}
}

func (l *FixedLayout) Appy(state *state.UIState, vm viewmodel.ViewModel, size terminal.Winsize) []text.Line {
	rows := min(l.maxRows, size.Rows)
	cols := min(l.maxCols, size.Cols)
	winsize := terminal.NewWinsize(rows, cols)
	return l.layout.Apply(state, vm, winsize)
}
