package text

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/app/pager"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/decorator/box"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/drain"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/focus"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/padding"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/widget/textarea/transformer"
	"github.com/Rafael24595/go-reacterm-core/engine/model/buffer/processor"
	"github.com/Rafael24595/go-reacterm-core/engine/model/hint"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

const NameInput = "text_input"

const input_limit = 20
const input_max_limit = 30

type TextInput struct {
	limit    winsize.Cols
	label    []text.Fragment
	textarea *TextArea
}

func NewInput() *TextInput {
	processor := processor.Limit(
		input_limit,
		processor.Inline,
	)

	area := NewArea().SetName(NameInput)
	area.buffer.Processor(processor)

	return &TextInput{
		limit:    input_limit,
		textarea: area,
	}
}

func (c *TextInput) SetName(name string) *TextInput {
	c.textarea.SetName(name)
	return c
}

func (c *TextInput) SetProcessor(limit winsize.Cols, process processor.Processor) *TextInput {
	assert.True(
		limit <= input_max_limit,
		"longer text fields should use the text area screen instead of the input one.",
	)

	c.limit = limit

	c.textarea.buffer.Processor(
		processor.Limit(limit, process),
	)

	return c
}

func (c *TextInput) SetLabel(label []text.Fragment) *TextInput {
	c.label = label
	return c
}

func (c *TextInput) WriteMode() *TextInput {
	c.textarea.WriteMode()
	return c
}

func (c *TextInput) ReadMode() *TextInput {
	c.textarea.ReadMode()
	return c
}

func (c *TextInput) EnableBlinking() *TextInput {
	c.textarea.EnableBlinking()
	return c
}

func (c *TextInput) DisableBlinking() *TextInput {
	c.textarea.DisableBlinking()
	return c
}

func (c *TextInput) AddText(text string) *TextInput {
	c.textarea.AddText(text)
	return c
}

func (c *TextInput) ToNode() screen.Node {
	return screen.NewBuilder().
		Name(c.textarea.reference).
		NameToStack().
		Definition(c.textarea.definition).
		Update(c.textarea.update).
		View(c.view).
		ToNode()
}

func (c *TextInput) view(stt state.UIState) viewmodel.ViewModel {
	vm := viewmodel.NewViewModel()

	_, textarea, needsPulse := c.textarea.viewSources()

	textarea.PushStep(
		transformer.BreakWord,
	)

	pipeline := c.makePipeline(
		textarea.ToUnit(),
	)

	box := box.New(pipeline).
		PaddingY(0).
		PaddingX(1).
		ToUnit()

	if len(c.label) != 0 {
		frags := append(c.label, *text.NewFragment(": "))
		vm.Kernel.Push(
			drain.UnitFromFragments(frags...),
		)
	}

	vm.Kernel.Push(box)

	vm.Behavior.NeedsPulse = needsPulse

	return *vm
}

func (c *TextInput) makePipeline(unit drawable.Unit) drawable.Unit {
	pageStep := pageTransformer()

	paddingStep := padding.ColsDrawTransformer(
		hint.Fixed(c.limit),
		style.Left,
	)

	return pipeline.New(unit).
		SetDrawStep(pageStep).
		PushDataSteps(paddingStep).
		ToUnit()
}

func limitRows(size winsize.Winsize) winsize.Winsize {
	return winsize.New(
		min(1, size.Rows),
		size.Cols,
	)
}

func pageTransformer() pipeline.DrawTransformer {
	engine := pager.EngineScroll()
	return func(winsize winsize.Winsize, unit drawable.Unit) ([]text.Line, bool) {
		transformer := focus.DrawTransformer(engine)
		return transformer(
			limitRows(winsize),
			unit,
		)
	}
}
