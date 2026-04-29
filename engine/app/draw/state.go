package draw

import (
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/work"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

type DrawState struct {
	Buffer  []text.Line
	HasNext bool
	Work    *work.Tracker
	Cursor  uint16
	Page    uint
	Focus   bool
}

func NewDrawStatus(ctx *DrawContext) *DrawState {
	return &DrawState{
		Buffer:  make([]text.Line, ctx.Size.Rows),
		HasNext: false,
		Work:    work.NewTracker(),
		Cursor:  0,
		Page:    0,
		Focus:   false,
	}
}

func (s *DrawState) MarkFocus(focus bool) *DrawState {
	s.Focus = s.Focus || focus
	return s
}

func (s *DrawState) SetAndNext(line text.Line) *DrawState {
	s.Buffer[s.Cursor] = line
	s.Cursor += 1
	return s
}

func (s *DrawState) Reset() {
	for i := range s.Buffer {
		s.Buffer[i] = text.Line{}
	}

	s.Cursor = 0
	s.Focus = false
}
