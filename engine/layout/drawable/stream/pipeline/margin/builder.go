package margin

import (
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

type builder struct {
	marginY   winsize.Rows
	positionY style.VerticalPosition
	marginX   winsize.Cols
	positionX style.HorizontalPosition
}

func NewBuilder() *builder {
	return &builder{
		marginY:   0,
		positionY: style.Middle,
		marginX:   0,
		positionX: style.Center,
	}
}

func (b *builder) Y(margin winsize.Rows, position ...style.VerticalPosition) *builder {
	b.marginY = margin

	if len(position) > 0 {
		b.positionY = position[0]
	}

	return b
}

func (b *builder) X(margin winsize.Cols, position ...style.HorizontalPosition) *builder {
	b.marginX = margin

	if len(position) > 0 {
		b.positionX = position[0]
	}

	return b
}

func (b *builder) Steps() (pipeline.DrawTransformer, []pipeline.DataTransformer) {
	draw := func(size winsize.Winsize, unit drawable.Unit) ([]text.Line, bool) {
		marginY := b.marginY * verticalFactor(b.positionY)
		marginX := b.marginX * horizontalFactor(b.positionX)

		fixedSize := winsize.New(
			size.Rows.Sub(marginY),
			size.Cols.Sub(marginX),
		)

		return unit.Drawable.Draw(fixedSize)
	}

	data := make([]pipeline.DataTransformer, 0, 2)

	if b.marginY > 0 {
		data = append(data,
			rowsDataTransformer(b.marginY, b.positionY),
		)
	}

	if b.marginX > 0 {
		data = append(data,
			colsDrawTransformer(b.marginX, b.positionX),
		)
	}

	return draw, data
}

func (b *builder) ToUnit(unit drawable.Unit) drawable.Unit {
	draw, data := b.Steps()
	return pipeline.New(unit).
		SetDrawStep(draw).
		PushDataSteps(data...).
		ToUnit()
}
