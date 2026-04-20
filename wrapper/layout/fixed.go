package wrapper_layout

import (
	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/app/viewmodel"
	"github.com/Rafael24595/go-terminal/engine/layout"
	"github.com/Rafael24595/go-terminal/engine/model/winsize"
	"github.com/Rafael24595/go-terminal/engine/render/text"
)

type FixedLayout struct {
	layout  layout.Layout
	maxRows winsize.Rows
	maxCols uint16
}

func NewFixed(layout layout.Layout, maxRows winsize.Rows, maxCols uint16) *FixedLayout {
	return &FixedLayout{
		layout:  layout,
		maxRows: maxRows,
		maxCols: maxCols,
	}
}

func (l *FixedLayout) Update(maxRows winsize.Rows, maxCols uint16) *FixedLayout {
	l.maxCols = maxCols
	l.maxRows = maxRows
	return l
}

func (l *FixedLayout) ToLayout() layout.Layout {
	return layout.Layout{
		Apply: l.Appy,
	}
}

func (l *FixedLayout) Appy(state *state.UIState, vm viewmodel.ViewModel, size winsize.Winsize) []text.Line {
	winsize := winsize.NewWinsize(
		min(l.maxRows, size.Rows),
		min(l.maxCols, size.Cols),
	)
	return l.layout.Apply(state, vm, winsize)
}
