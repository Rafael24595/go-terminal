package padding

import (
	"github.com/Rafael24595/go-reacterm-core/engine/config/padding/cols"
	"github.com/Rafael24595/go-reacterm-core/engine/config/padding/rows"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline"
	"github.com/Rafael24595/go-reacterm-core/engine/model/hint"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
)

type Builder struct {
	hintY    *hint.Size[winsize.Rows]
	optionsY []rows.Option
	hintX    *hint.Size[winsize.Cols]
	optionsX []cols.Option
}

func NewBuilder() *Builder {
	return &Builder{
		hintY:    nil,
		optionsY: make([]rows.Option, 0),
		hintX:    nil,
		optionsX: make([]cols.Option, 0),
	}
}

func (b *Builder) Y(hint hint.Size[winsize.Rows], opts ...rows.Option) *Builder {
	b.hintY = &hint
	b.optionsY = append(b.optionsY, opts...)
	return b
}

func (b *Builder) X(hint hint.Size[winsize.Cols], opts ...cols.Option) *Builder {
	b.hintX = &hint
	b.optionsX = append(b.optionsX, opts...)
	return b
}

func (b *Builder) Steps() []pipeline.DataTransformer {
	data := make([]pipeline.DataTransformer, 0, 2)

	if b.hintY != nil {
		data = append(data,
			Rows(*b.hintY, b.optionsY...),
		)
	}

	if b.hintX != nil {
		data = append(data,
			Cols(*b.hintX, b.optionsX...),
		)
	}

	return data
}

func (b *Builder) ToUnit(unit drawable.Unit) drawable.Unit {
	return pipeline.New(unit).
		PushDataSteps(b.Steps()...).
		ToUnit()
}
