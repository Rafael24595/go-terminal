package margin

import (
	"github.com/Rafael24595/go-reacterm-core/engine/config/padding/cols"
	"github.com/Rafael24595/go-reacterm-core/engine/config/padding/rows"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

type Builder struct {
	marginY  winsize.Rows
	optionsY []rows.Option
	marginX  winsize.Cols
	optionsX []cols.Option
}

func NewBuilder() *Builder {
	return &Builder{
		marginY:  0,
		optionsY: make([]rows.Option, 0),
		marginX:  0,
		optionsX: make([]cols.Option, 0),
	}
}

func (b *Builder) Y(margin winsize.Rows, opts ...rows.Option) *Builder {
	b.marginY = margin
	b.optionsY = append(b.optionsY, opts...)
	return b
}

func (b *Builder) X(margin winsize.Cols, opts ...cols.Option) *Builder {
	b.marginX = margin
	b.optionsX = append(b.optionsX, opts...)
	return b
}

func (b *Builder) Steps() (pipeline.DrawTransformer, []pipeline.DataTransformer) {
	draw := func(size winsize.Winsize, unit drawable.Unit) ([]text.Line, bool) {
		cfgY := rows.ResolveConfig(b.optionsY...)
		cfgX := cols.ResolveConfig(b.optionsX...)

		marginY := b.marginY * verticalFactor(cfgY.Position)
		marginX := b.marginX * horizontalFactor(cfgX.Position)

		fixedSize := winsize.New(
			size.Rows.Sub(marginY),
			size.Cols.Sub(marginX),
		)

		return unit.Drawable.Draw(fixedSize)
	}

	data := make([]pipeline.DataTransformer, 0, 2)

	if b.marginY > 0 {
		data = append(data,
			rowsDataTransformer(b.marginY, b.optionsY...),
		)
	}

	if b.marginX > 0 {
		data = append(data,
			colsDrawTransformer(b.marginX, b.optionsX...),
		)
	}

	return draw, data
}

func (b *Builder) ToUnit(unit drawable.Unit) drawable.Unit {
	draw, data := b.Steps()
	return pipeline.New(unit).
		SetDrawStep(draw).
		PushDataSteps(data...).
		ToUnit()
}
