package draw

import (
	"github.com/Rafael24595/go-terminal/engine/commons/structure/work"
	"github.com/Rafael24595/go-terminal/engine/render/text"
)

type DrawState struct {
	Buffer []text.Line
	Work   *work.Tracker
	Cursor uint16
	Page   uint
	Focus  bool
}

func NewDrawStatus(ctx *DrawContext) *DrawState {
	return &DrawState{
		Buffer: make([]text.Line, ctx.Size.Rows),
		Work:   work.NewTracker(),
		Cursor: 0,
		Page:   0,
		Focus:  false,
	}
}
