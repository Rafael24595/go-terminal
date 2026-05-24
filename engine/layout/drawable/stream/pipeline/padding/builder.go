package padding

import (
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/utils/padding"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/utils/padding/options"
	"github.com/Rafael24595/go-reacterm-core/engine/model/hint"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
)

type Builder struct {
	hintY     *hint.Size[winsize.Rows]
	optionsY  []options.RowsOption
	hintX     *hint.Size[winsize.Cols]
	defaultX  string
	positionX style.HorizontalPosition
}

func NewBuilder() *Builder {
	return &Builder{
		hintY:     nil,
		optionsY:  make([]options.RowsOption, 0),
		hintX:     nil,
		defaultX:  padding.DefaultColFrag,
		positionX: style.Center,
	}
}

func (b *Builder) Y(hint hint.Size[winsize.Rows], opts ...options.RowsOption) *Builder {
	b.hintY = &hint
	b.optionsY = append(b.optionsY, opts...)
	return b
}

func (b *Builder) X(hint hint.Size[winsize.Cols], position ...style.HorizontalPosition) *Builder {
	return b.XWithDefault(hint, padding.DefaultColFrag, position...)
}

func (b *Builder) XWithDefault(hint hint.Size[winsize.Cols], frag string, position ...style.HorizontalPosition) *Builder {
	b.hintX = &hint
	b.defaultX = frag

	if len(position) > 0 {
		b.positionX = position[0]
	}

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
			ColsWithDefault(*b.hintX, b.defaultX, b.positionX),
		)
	}

	return data
}

func (b *Builder) ToUnit(unit drawable.Unit) drawable.Unit {
	return pipeline.New(unit).
		PushDataSteps(b.Steps()...).
		ToUnit()
}
