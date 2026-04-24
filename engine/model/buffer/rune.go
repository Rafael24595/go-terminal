package buffer

import (
	"github.com/Rafael24595/go-reacterm-core/engine/helper/runes"
	"github.com/Rafael24595/go-reacterm-core/engine/helper/text"
	"github.com/Rafael24595/go-reacterm-core/engine/model/delta"
)

type RuneBuffer struct {
	buffer      []rune
	facade      []rune
	transformer text.TextTransformer
	handler     RuneHandler
}

func NewRuneBuffer() *RuneBuffer {
	return &RuneBuffer{
		buffer:      make([]rune, 0),
		facade:      make([]rune, 0),
		transformer: text.VoidTextTransformer,
		handler:     voidRuneHandler,
	}
}

func (b *RuneBuffer) Transformer(transformer text.TextTransformer) *RuneBuffer {
	b.transformer = transformer
	return b
}

func (b *RuneBuffer) Handler(handler RuneHandler) *RuneBuffer {
	b.handler = handler
	return b
}

func (b *RuneBuffer) Size() uint {
	return uint(len(b.buffer))
}

func (b *RuneBuffer) Empty() bool {
	return len(b.buffer) == 0
}

func (b *RuneBuffer) Buffer() []rune {
	return b.buffer
}

func (b *RuneBuffer) Facade() []rune {
	return b.facade
}

func (b *RuneBuffer) Range(start uint, end uint) []rune {
	if end < start {
		return make([]rune, 0)
	}

	return b.buffer[start:end]
}

func (b *RuneBuffer) Append(rns []rune) *RuneBuffer {
	b.Replace(rns, b.Size(), b.Size())
	return b
}

func (b *RuneBuffer) TransformAndReplace(rns []rune, start uint, end uint) ([]rune, []rune) {
	if end < start {
		zero := make([]rune, 0)
		return zero, zero
	}

	insert := b.transformer.Apply(rns, start, end, b.buffer)
	return b.applyChange(insert, start, end)
}

func (b *RuneBuffer) Replace(rns []rune, start uint, end uint) ([]rune, []rune) {
	if end < start {
		zero := make([]rune, 0)
		return zero, zero
	}

	return b.applyChange(rns, start, end)
}

func (b *RuneBuffer) Delete(start uint, end uint) []rune {
	if end < start {
		return make([]rune, 0)
	}

	rns := make([]rune, 0)
	_, deleted := b.Replace(rns, start, end)
	return deleted
}

func (b *RuneBuffer) applyChange(insert []rune, start, end uint) ([]rune, []rune) {
	if end > uint(len(b.buffer)) {
		end = uint(len(b.buffer))
	}

	deleted := b.Range(start, end)

	rawBuffer := runes.AppendRange(b.buffer, insert, start, end)
	newBuffer, newFacade := b.handler(rawBuffer)

	insertSize := len(newBuffer) - (len(b.buffer) - len(deleted))

	fixedInsert := make([]rune, 0)
	if insertSize > 0 {
		endInsert := min(start+uint(insertSize), uint(len(newBuffer)))
		fixedInsert = newBuffer[start:endInsert]
	}

	b.buffer = newBuffer
	b.facade = newFacade

	return fixedInsert, deleted
}

func (b *RuneBuffer) ApplyDelta(d *delta.Delta) *RuneBuffer {
	newBuffer := delta.Apply(b.buffer, d)
	buffer, facade := b.handler(newBuffer)

	b.buffer = buffer
	b.facade = facade

	return b
}
