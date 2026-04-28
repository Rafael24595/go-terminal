package text

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"
	
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/decorator/box"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/spatial/position"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/block"
	"github.com/Rafael24595/go-reacterm-core/engine/model/buffer"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"

	text_transformer "github.com/Rafael24595/go-reacterm-core/engine/helper/text"
)

const input_name = "TextInput"
const input_limit = 20

const input_max_limit = 30

type TextInput struct {
	limit    uint64
	label    []text.Fragment
	textarea *TextArea
}

func NewInput() *TextInput {
	handler := buffer.NewLimitedRuneHandler(input_limit, buffer.String)

	area := NewArea().SetName(input_name)
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

func (c *TextInput) SetType(limit uint64, input buffer.InputType) *TextInput {
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

	vm.Kernel.Unshift(
		c.makeDrawables(vm)...,
	)

	return vm
}

func (c *TextInput) makeDrawables(vm viewmodel.ViewModel) []drawable.Drawable {
	dws := make([]drawable.Drawable, 0, 2)

	code := c.textarea.mainDrawableCode()
	dw, ok := vm.Kernel.Take(code)
	if !ok {
		return dws
	}

	inp := box.NewBoxDrawable(dw).
		PaddingY(0).
		PaddingX(1).
		TextAlign(style.Left).
		MinSize(uint(c.limit)).
		ToDrawable()

	pst := position.NewPositionDrawable(inp).
		PositionY(style.Top).
		PositionX(style.Left)

	if len(c.label) == 0 {
		dws = append(dws, pst.ToDrawable())
		return dws
	}

	//TODO: Parametrize.
	pstd := pst.MarginX(0).ToDrawable()

	frags := append(c.label, *text.NewFragment(": "))
	lbl := block.BlockDrawableFromLines(
		*text.LineFromFragments(frags...),
	)

	dws = append(dws, lbl)
	dws = append(dws, pstd)

	return dws
}
