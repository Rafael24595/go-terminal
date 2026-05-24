package margin

import (
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/utils/padding"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/utils/padding/options"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

type Builder struct {
	marginY   winsize.Rows
	optionsY  []options.RowsOption
	marginX   winsize.Cols
	defaultX  string
	positionX style.HorizontalPosition
}

func NewBuilder() *Builder {
	return &Builder{
		marginY:   0,
		optionsY:  make([]options.RowsOption, 0),
		marginX:   0,
		defaultX:  padding.DefaultColFrag,
		positionX: style.Center,
	}
}

func (b *Builder) Y(margin winsize.Rows, opts ...options.RowsOption) *Builder {
	b.marginY = margin
	b.optionsY = append(b.optionsY, opts...)
	return b
}

func (b *Builder) X(margin winsize.Cols, position ...style.HorizontalPosition) *Builder {
	return b.XWithDefault(margin, padding.DefaultColFrag, position...)
}

func (b *Builder) XWithDefault(margin winsize.Cols, frag string, position ...style.HorizontalPosition) *Builder {
	b.marginX = margin
	b.defaultX = frag

	if len(position) > 0 {
		b.positionX = position[0]
	}

	return b
}

func (b *Builder) Steps() (pipeline.DrawTransformer, []pipeline.DataTransformer) {
	draw := func(size winsize.Winsize, unit drawable.Unit) ([]text.Line, bool) {
		cfg := options.ResolveRowsConfig(b.optionsY...)

		marginY := b.marginY * verticalFactor(cfg.Position)
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
			rowsDataTransformer(b.marginY, b.optionsY...),
		)
	}

	if b.marginX > 0 {
		data = append(data,
			colsDrawTransformer(b.marginX, b.defaultX, b.positionX),
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
