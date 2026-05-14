package text

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/decorator/box"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/spatial/position"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/drain"
	"github.com/Rafael24595/go-reacterm-core/engine/model/buffer"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"

	text_transformer "github.com/Rafael24595/go-reacterm-core/engine/helper/text"
)

const NameInput = "text_input"

const input_limit = 200
const input_max_limit = 30

type TextInput struct {
	limit    winsize.Cols
	label    []text.Fragment
	textarea *TextArea
}

func NewInput() *TextInput {
	handler := buffer.NewLimitedRuneHandler(input_limit, buffer.String)

	area := NewArea().SetName(NameInput)
	area.buffer.Handler(handler)
	area.buffer.Transformer(text_transformer.VoidTextTransformer)

	return &TextInput{
		limit:    input_limit,
		textarea: area,
	}
}

func (c *TextInput) SetName(name string) *TextInput {
	c.textarea.SetName(name)
	return c
}

func (c *TextInput) SetType(limit winsize.Cols, input buffer.InputType) *TextInput {
	assert.True(
		limit <= input_max_limit,
		"longer text fields should use the text area screen instead of the input one.",
	)

	c.limit = limit

	handler := buffer.NewLimitedRuneHandler(limit, input)
	c.textarea.buffer.Handler(handler)

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
	vm := c.textarea.view(stt)

	vm.Kernel.Unshift(
		c.makeDrawables(vm)...,
	)

	return vm
}

func (c *TextInput) makeDrawables(vm viewmodel.ViewModel) []drawable.Drawable {
	drawables := make([]drawable.Drawable, 0, 2)

	code := c.textarea.mainDrawableCode()
	drawable, ok := vm.Kernel.Take(code)
	if !ok {
		return drawables
	}

	input := box.New(drawable).
		PaddingY(0).
		PaddingX(1).
		TextAlign(style.Left).
		MinSize(c.limit).
		ToDrawable()

	position := position.New(input).
		PositionY(style.Top).
		PositionX(style.Left)

	if len(c.label) == 0 {
		drawables = append(drawables, position.ToDrawable())
		return drawables
	}

	frags := append(c.label, *text.NewFragment(": "))
	
	drawables = append(drawables, 
		drain.DrawableFromFragments(frags...),
	)
	
	//TODO: Parametrize.
	drawables = append(drawables, 
		position.MarginX(0).ToDrawable(),
	)

	return drawables
}
