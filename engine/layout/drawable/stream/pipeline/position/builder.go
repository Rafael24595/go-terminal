package position

import (
	"github.com/Rafael24595/go-reacterm-core/engine/config/padding/cols"
	"github.com/Rafael24595/go-reacterm-core/engine/config/padding/rows"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/margin"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/padding"
	"github.com/Rafael24595/go-reacterm-core/engine/model/hint"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
)

type Builder struct {
	padding *padding.Builder
	margin  *margin.Builder
}

func New() *Builder {
	return &Builder{
		padding: padding.NewBuilder(),
		margin:  margin.NewBuilder(),
	}
}

func (b *Builder) MarginY(margin winsize.Rows, opts ...rows.Option) *Builder {
	b.margin.Y(margin, opts...)
	return b
}

func (b *Builder) MarginX(margin winsize.Cols, opts ...cols.Option) *Builder {
	b.margin.X(margin, opts...)
	return b
}

func (b *Builder) PaddingY(hint hint.Size[winsize.Rows], opts ...rows.Option) *Builder {
	b.padding.Y(hint, opts...)
	return b
}

func (b *Builder) PaddingX(hint hint.Size[winsize.Cols], opts ...cols.Option) *Builder {
	b.padding.X(hint, opts...)
	return b
}

func (b *Builder) ToUnit(unit drawable.Unit) drawable.Unit {
	marginDraw, marginDatas := b.margin.Steps()
	paddingDatas := b.padding.Steps()

	return pipeline.New(unit).
		SetDrawStep(marginDraw).
		PushDataSteps(marginDatas...).
		PushDataSteps(paddingDatas...).
		ToUnit()
}
