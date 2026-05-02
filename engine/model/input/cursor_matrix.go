package input

import (
	"github.com/Rafael24595/go-reacterm-core/engine/helper/math"
)

type MatrixCursor struct {
	Row  uint16
	Col  uint16
	Show bool
}

func NewMatrixCursor(row, col uint16, show bool) *MatrixCursor {
	return &MatrixCursor{
		Row:  row,
		Col:  col,
		Show: show,
	}
}

func (c *MatrixCursor) IncRow(limit uint16) *MatrixCursor {
	c.Row = min(limit, c.Row+1)
	return c
}

func (c *MatrixCursor) DecRow() *MatrixCursor {
	c.Row = math.SubClampZero(c.Row, 1)
	return c
}

func (c *MatrixCursor) IncCol(limit uint16) *MatrixCursor {
	c.Col = min(limit, c.Col+1)
	return c
}

func (c *MatrixCursor) DecCol() *MatrixCursor {
	c.Col = math.SubClampZero(c.Col, 1)
	return c
}
