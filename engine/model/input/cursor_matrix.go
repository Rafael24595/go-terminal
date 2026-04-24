package input

import "github.com/Rafael24595/go-reacterm-core/engine/helper/math"

type MatrixCursor struct {
	Row  uint32
	Col  uint32
	Show bool
}

func NewMatrixCursor(row uint32, col uint32, show bool) *MatrixCursor {
	return &MatrixCursor{
		Row:  row,
		Col:  col,
		Show: show,
	}
}

func (c *MatrixCursor) IncRow(len uint32) *MatrixCursor {
	c.Row = min(len, c.Row+1)
	return c
}

func (c *MatrixCursor) DecRow() *MatrixCursor {
	c.Row = math.SubClampZero(c.Row, 1)
	return c
}

func (c *MatrixCursor) IncCol(len uint32) *MatrixCursor {
	c.Col = min(len, c.Col+1)
	return c
}

func (c *MatrixCursor) DecCol() *MatrixCursor {
	c.Col = math.SubClampZero(c.Col, 1)
	return c
}
