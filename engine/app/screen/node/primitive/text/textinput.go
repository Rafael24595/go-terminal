package text

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/app/pager"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/config/padding/cols"
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

func (n *TextInput) SetName(name string) *TextInput {
	n.textarea.SetName(name)
	return n
}

func (n *TextInput) SetProcessor(limit winsize.Cols, process processor.Processor) *TextInput {
	assert.True(
		limit <= input_max_limit,
		"longer text fields should use the text area screen instead of the input one.",
	)

	n.limit = limit

	n.textarea.buffer.Processor(
		processor.Limit(limit, process),
	)

	return n
}

func (n *TextInput) SetLabel(label []text.Fragment) *TextInput {
	n.label = label
	return n
}

func (n *TextInput) WriteMode() *TextInput {
	n.textarea.WriteMode()
	return n
}

func (n *TextInput) ReadMode() *TextInput {
	n.textarea.ReadMode()
	return n
}

func (n *TextInput) EnableBlinking() *TextInput {
	n.textarea.EnableBlinking()
	return n
}

func (n *TextInput) DisableBlinking() *TextInput {
	n.textarea.DisableBlinking()
	return n
}

func (n *TextInput) AddText(text string) *TextInput {
	n.textarea.AddText(text)
	return n
}

func (n *TextInput) ToNode() screen.Node {
	return screen.NewBuilder().
		Name(n.textarea.reference).
		NameToStack().
		Keys(n.textarea.keys).
		Tick(n.textarea.tick).
		View(n.view).
		ToNode()
}

func (n *TextInput) view(uiState state.UIState) viewmodel.ViewModel {
	vm := viewmodel.New()

	_, textarea, needsPulse := n.textarea.viewSources()

	textarea.PushStep(
		transformer.BreakWord,
	)

	pipeline := n.makePipeline(
		textarea.ToUnit(),
	)

	box := box.New(pipeline).
		PaddingY(0).
		PaddingX(1).
		ToUnit()

	if len(n.label) != 0 {
		frags := append(n.label, *text.NewFragment(": "))
		vm.Kernel.Push(
			drain.UnitFromFragments(frags...),
		)
	}

	vm.Kernel.Push(box)

	vm.Behavior.NeedsPulse = needsPulse

	return *vm
}

func (n *TextInput) makePipeline(unit drawable.Unit) drawable.Unit {
	pageStep := pageTransformer()

	paddingStep := padding.Cols(
		hint.Fixed(n.limit),
		cols.WithPosition(style.Left),
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
