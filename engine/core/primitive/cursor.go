package primitive

import (
	"time"

	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/helper/math"
)

const blink_ms = 750

type clock func() int64

type Cursor struct {
	clock  clock
	blink  bool
	status bool
	time   int64
	caret  uint
	anchor uint
}

func NewCursor(blink bool) *Cursor {
	return &Cursor{
		clock:  time.Now().UnixMilli,
		blink:  blink,
		status: true,
		time:   0,
		caret:  0,
		anchor: 0,
	}
}

func (c *Cursor) EnableBlinking() *Cursor {
	c.blink = true
	return c
}

func (c *Cursor) DisableBlinking() *Cursor {
	c.blink = false
	return c
}

func (c *Cursor) Caret() uint {
	return c.caret
}

func (c *Cursor) Anchor() uint {
	return c.anchor
}

func (c *Cursor) SelectStart() uint {
	if c.anchor < c.caret {
		return c.anchor
	}
	return c.caret
}

func (c *Cursor) SelectEnd() uint {
	if c.anchor < c.caret {
		return c.caret
	}
	return c.anchor
}

func (c *Cursor) MoveCaretTo(buff []rune, caret uint) {
	min := uint(1)
	len := uint(len(buff))

	if len == 0 {
		min = 0
	}

	c.caret = math.Clamp(caret, min, len)
	c.anchor = c.caret

	c.status = true
	c.time = c.clock()
}

func (c *Cursor) MoveSelectTo(buff []rune, caret, anchor uint) {
	min := uint(1)
	len := uint(len(buff))

	if len == 0 {
		min = 0
	}

	c.caret = math.Clamp(caret, min, len)
	c.anchor = math.Clamp(anchor, min, len)

	c.status = true
	c.time = c.clock()
}

func (c *Cursor) BlinkStyle() core.Style {
	if !c.blink || c.caret != c.anchor {
		return core.Select
	}

	style := core.None
	if c.status {
		style = core.Select
	}

	now := c.clock()
	if now-c.time >= blink_ms {
		c.time = now
		c.status = !c.status
	}

	return style
}
