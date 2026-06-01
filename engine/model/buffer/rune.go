package buffer

import (
	"github.com/Rafael24595/go-reacterm-core/engine/helper/runes"
	"github.com/Rafael24595/go-reacterm-core/engine/model/buffer/processor"
	"github.com/Rafael24595/go-reacterm-core/engine/model/buffer/rule"
	"github.com/Rafael24595/go-reacterm-core/engine/model/delta"
	"github.com/Rafael24595/go-reacterm-core/engine/model/offset"
)

type RuneBuffer struct {
	buffer    []rune
	facade    []rune
	rules     []rule.Rule
	processor processor.Processor
}

func NewRuneBuffer() *RuneBuffer {
	return &RuneBuffer{
		buffer:    make([]rune, 0),
		facade:    make([]rune, 0),
		rules:     make([]rule.Rule, 0),
		processor: processor.Identity,
	}
}

func (b *RuneBuffer) PushRules(rules ...rule.Rule) *RuneBuffer {
	b.rules = append(b.rules, rules...)
	return b
}

func (b *RuneBuffer) Processor(processor processor.Processor) *RuneBuffer {
	b.processor = processor
	return b
}

func (b *RuneBuffer) Size() offset.Offset {
	return offset.Offset(len(b.buffer))
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

func (b *RuneBuffer) Range(start offset.Offset, end offset.Offset) []rune {
	if end < start {
		return make([]rune, 0)
	}

	return b.buffer[start:end]
}

func (b *RuneBuffer) Append(rns []rune) *RuneBuffer {
	b.Replace(rns, b.Size(), b.Size())
	return b
}

func (b *RuneBuffer) Clean() *RuneBuffer {
	b.buffer = make([]rune, 0)
	b.facade = make([]rune, 0)
	return b
}

func (b *RuneBuffer) Replace(rns []rune, start offset.Offset, end offset.Offset) ([]rune, []rune) {
	if end < start {
		zero := make([]rune, 0)
		return zero, zero
	}

	return b.commitReplace(rns, start, end)
}

func (b *RuneBuffer) ReplaceWithRules(buffer []rune, start offset.Offset, end offset.Offset) ([]rune, []rune) {
	if end < start {
		zero := make([]rune, 0)
		return zero, zero
	}

	insert := b.applyRules(buffer, start, end, b.buffer)
	return b.commitReplace(insert, start, end)
}

func (b *RuneBuffer) applyRules(text []rune, start, end offset.Offset, buff []rune) []rune {
	for _, rule := range b.rules {
		if text, ok := rule(text, start, end, buff); ok {
			return text
		}
	}
	return text
}

func (b *RuneBuffer) Delete(start offset.Offset, end offset.Offset) []rune {
	if end < start {
		return make([]rune, 0)
	}

	rns := make([]rune, 0)
	_, deleted := b.Replace(rns, start, end)
	return deleted
}

func (b *RuneBuffer) commitReplace(insert []rune, start, end offset.Offset) ([]rune, []rune) {
	end = min(end, offset.Offset(len(b.buffer)))

	deleted := b.Range(start, end)

	rawBuffer := runes.AppendRange(b.buffer, insert, start, end)
	newBuffer, newFacade := b.processor(rawBuffer)

	newBufferLen := offset.Offset(len(newBuffer))

	insertSize := offset.Offset(
		len(newBuffer) - (len(b.buffer) - len(deleted)),
	)

	fixedInsert := make([]rune, 0)
	if insertSize > 0 {
		endInsert := min(start+insertSize, newBufferLen)
		fixedInsert = newBuffer[start:endInsert]
	}

	b.buffer = newBuffer
	b.facade = newFacade

	return fixedInsert, deleted
}

func (b *RuneBuffer) ApplyDelta(d *delta.Delta) *RuneBuffer {
	newBuffer := delta.Apply(b.buffer, d)
	buffer, facade := b.processor(newBuffer)

	b.buffer = buffer
	b.facade = facade

	return b
}
