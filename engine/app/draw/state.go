package draw

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"
	
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
	if s.IsFull() {
		assert.Unreachable("buffer overflow")
		return s
	}

	s.Buffer[s.Cursor] = line
	s.Cursor += 1
	return s
}

func (s *DrawState) IsFull() bool {
	return s.Cursor == uint16(len(s.Buffer))
}

func (s *DrawState) Reset() {
	for i := range s.Buffer {
		s.Buffer[i] = text.Line{}
	}

	s.Cursor = 0
	s.Focus = false
}
