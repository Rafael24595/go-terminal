package table

import "github.com/Rafael24595/go-terminal/engine/helper/math"

type Cursor struct {
	Row  uint32
	Col  uint32
	Show bool
}

func NewCursor(row uint32, col uint32, show bool) *Cursor {
	return &Cursor{
		Row:  row,
		Col:  col,
		Show: show,
	}
}

func (c *Cursor) IncRow(len uint32) *Cursor {
	c.Row = min(len, c.Row+1)
	return c
}

func (c *Cursor) DecRow() *Cursor {
	c.Row = math.SubClampZero(c.Row, 1)
	return c
}

func (c *Cursor) IncCol(len uint32) *Cursor {
	c.Col = min(len, c.Col+1)
	return c
}

func (c *Cursor) DecCol() *Cursor {
	c.Col = math.SubClampZero(c.Col, 1)
	return c
}
