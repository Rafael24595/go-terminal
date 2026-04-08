package viewmodel

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-terminal/engine/commons/structure/set"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable"
	"github.com/Rafael24595/go-terminal/engine/render/marker"
	"github.com/Rafael24595/go-terminal/engine/render/text"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

const NameInputLineDrawable = "InputLineDrawable"

type InputLine struct {
	Prompt string
	Value  drawable.Drawable
}

func NewInputLine(drawable drawable.Drawable) *InputLine {
	return &InputLine{
		Prompt: marker.DefaultInputLinePrompt,
		Value:  drawable,
	}
}

func (i *InputLine) ToDrawable() drawable.Drawable {
	return toDrawable(i)
}

type inputLineDrawable struct {
	loaded bool
	status bool
	prompt string
	input  drawable.Drawable
}

func toDrawable(input *InputLine) drawable.Drawable {
	drw := inputLineDrawable{
		loaded: false,
		status: true,
		prompt: input.Prompt,
		input:  input.Value,
	}
	return drawable.Drawable{
		Name: NameInputLineDrawable,
		Code: "",
		Tags: make(set.Set[string]),
		Init: drw.init,
		Wipe: drw.wipe,
		Draw: drw.draw,
	}
}

func (d *inputLineDrawable) init() {
	d.loaded = true

	d.input.Init()
}

func (d *inputLineDrawable) wipe() {}

func (d *inputLineDrawable) draw(size terminal.Winsize) ([]text.Line, bool) {
	assert.True(d.loaded, "the drawable should be initialized before draw")

	lines := make([]text.Line, 0)

	if !d.status {
		return lines, false
	}

	content := true
	for content {
		result, status := d.input.Draw(size)
		content = status
		lines = append(lines, result...)
	}

	if len(lines) == 0 {
		line := text.NewLine(d.prompt)
		return []text.Line{line}, false
	}

	prompt := text.FragmentsFromString(d.prompt + marker.DefaultPaddingText)
	lines[0].Text = append(prompt, lines[0].Text...)

	return append([]text.Line{text.EmptyLine()}, lines...), false
}
