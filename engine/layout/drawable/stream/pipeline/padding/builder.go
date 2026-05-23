package padding

import (
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline"
	"github.com/Rafael24595/go-reacterm-core/engine/model/hint"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
)

type builder struct {
	hintY     *hint.Size[winsize.Rows]
	positionY style.VerticalPosition
	hintX     *hint.Size[winsize.Cols]
	positionX style.HorizontalPosition
}

func NewBuilder() *builder {
	return &builder{
		hintY:     nil,
		positionY: style.Middle,
		hintX:     nil,
		positionX: style.Center,
	}
}

func (b *builder) Y(hint hint.Size[winsize.Rows], position ...style.VerticalPosition) *builder {
	b.hintY = &hint

	if len(position) > 0 {
		b.positionY = position[0]
	}

	return b
}

func (b *builder) X(hint hint.Size[winsize.Cols], position ...style.HorizontalPosition) *builder {
	b.hintX = &hint

	if len(position) > 0 {
		b.positionX = position[0]
	}

	return b
}

func (b *builder) Steps() []pipeline.DataTransformer {
	data := make([]pipeline.DataTransformer, 0, 2)

	if b.hintY != nil {
		data = append(data,
			Rows(*b.hintY, b.positionY),
		)
	}

	if b.hintX != nil {
		data = append(data,
			Cols(*b.hintX, b.positionX),
		)
	}

	return data
}

func (b *builder) ToUnit(unit drawable.Unit) drawable.Unit {
	return pipeline.New(unit).
		PushDataSteps(b.Steps()...).
		ToUnit()
}
