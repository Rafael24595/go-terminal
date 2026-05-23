package padding

import (
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/utils/padding"
	"github.com/Rafael24595/go-reacterm-core/engine/model/hint"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

type Builder struct {
	hintY     *hint.Size[winsize.Rows]
	defaultY  text.Fragment
	positionY style.VerticalPosition
	hintX     *hint.Size[winsize.Cols]
	defaultX  string
	positionX style.HorizontalPosition
}

func NewBuilder() *Builder {
	return &Builder{
		hintY:     nil,
		defaultY:  *padding.DefaultRowFrag(),
		positionY: style.Middle,
		hintX:     nil,
		defaultX:  padding.DefaultColFrag,
		positionX: style.Center,
	}
}

func (b *Builder) Y(hint hint.Size[winsize.Rows], position ...style.VerticalPosition) *Builder {
	return b.YWithDefault(hint, *padding.DefaultRowFrag(), position...)
}

func (b *Builder) YWithDefault(hint hint.Size[winsize.Rows], frag text.Fragment, position ...style.VerticalPosition) *Builder {
	b.hintY = &hint
	b.defaultY = frag

	if len(position) > 0 {
		b.positionY = position[0]
	}

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
			RowsWithDefault(*b.hintY, b.defaultY, b.positionY),
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
