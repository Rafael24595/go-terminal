package input

import (
	"github.com/Rafael24595/go-reacterm-core/engine/helper/math"
	"github.com/Rafael24595/go-reacterm-core/engine/platform/clock"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
)

const blink_ms = 750

type TextCursor struct {
	clock  clock.Clock
	blink  bool
	status bool
	time   int64
	caret  uint
	anchor uint
}

func NewTextCursor(blink bool) *TextCursor {
	return &TextCursor{
		clock:  clock.UnixMilliClock,
		blink:  blink,
		status: true,
		time:   0,
		caret:  0,
		anchor: 0,
	}
}

func (c *TextCursor) IsBlinking() bool {
	return c.blink
}

func (c *TextCursor) EnableBlinking() *TextCursor {
	c.blink = true
	return c
}

func (c *TextCursor) DisableBlinking() *TextCursor {
	c.blink = false
	return c
}

func (c *TextCursor) Caret() uint {
	return c.caret
}

func (c *TextCursor) Anchor() uint {
	return c.anchor
}

func (c *TextCursor) SelectStart() uint {
	if c.anchor < c.caret {
		return c.anchor
	}
	return c.caret
}

func (c *TextCursor) SelectEnd() uint {
	if c.anchor < c.caret {
		return c.caret
	}
	return c.anchor
}

func (c *TextCursor) MoveCaretTo(buff []rune, caret uint) {
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

func (c *TextCursor) MoveSelectTo(buff []rune, caret, anchor uint) {
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

func (c *TextCursor) BlinkStyle() style.Atom {
	if !c.blink || c.caret != c.anchor {
		return style.AtmSelect
	}

	styl := style.AtmNone
	if c.status {
		styl = style.AtmSelect
	}

	now := c.clock()
	if now-c.time >= blink_ms {
		c.time = now
		c.status = !c.status
	}

	return styl
}
