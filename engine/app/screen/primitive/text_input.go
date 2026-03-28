package primitive

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"
	"github.com/Rafael24595/go-terminal/engine/app/screen"
	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/app/viewmodel"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable/box"
	"github.com/Rafael24595/go-terminal/engine/model/buffer"
	"github.com/Rafael24595/go-terminal/engine/render/style"

	text_transformer "github.com/Rafael24595/go-terminal/engine/helper/text"
)

const default_text_input_name = "TextInput"
const default_text_input_limit = 20

const text_input_max_limit = 30

type TextInput struct {
	limit    uint64
	textarea *TextArea
}

func NewTextInput() *TextInput {
	handler := buffer.NewLimitedRuneHandler(default_text_input_limit, buffer.String)

	textarea := NewTextArea().SetName(default_text_input_name)
	textarea.buffer.Handler(handler)
	textarea.buffer.Transformer(text_transformer.VoidTextTransformer)

	return &TextInput{
		limit:    default_text_input_limit,
		textarea: textarea,
	}
}

func (c *TextInput) SetName(name string) *TextInput {
	c.textarea.SetName(name)
	return c
}

func (c *TextInput) SetType(limit uint64, input buffer.InputType) *TextInput {
	assert.True(
		limit <= text_input_max_limit,
		"longer text fields should use the text area screen instead of the input one.",
	)

	c.limit = limit

	handler := buffer.NewLimitedRuneHandler(limit, input)
	c.textarea.buffer.Handler(handler)

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

func (c *TextInput) ToScreen() screen.Screen {
	screen := screen.Screen{
		Definition: c.textarea.definition,
		Update:     c.textarea.update,
		View:       c.view,
	}

	return screen.SetName(c.textarea.reference).
		StackFromName()
}

func (c *TextInput) view(stt state.UIState) viewmodel.ViewModel {
	vm := c.textarea.view(stt)

	code := c.textarea.mainDrawableCode()

	dw, ok := vm.Kernel.Take(code)
	if !ok {
		return vm
	}

	bx := box.NewBoxDrawable(dw).
		PaddingY(0).
		PaddingX(1).
		TextAlign(style.Left).
		MinSize(uint(c.limit))

	vm.Kernel.Unshift(bx.ToDrawable())

	return vm
}
