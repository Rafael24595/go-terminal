package core

import (
	"github.com/Rafael24595/go-terminal/engine/core/assert"
	"github.com/Rafael24595/go-terminal/engine/core/drawable"
	"github.com/Rafael24595/go-terminal/engine/core/marker"
	"github.com/Rafael24595/go-terminal/engine/core/text"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

const DefaultPrompt = ">"

type InputLine struct {
	Prompt string
	Value  drawable.Drawable
}

func NewInputLine(drawable drawable.Drawable) *InputLine {
	return &InputLine{
		Prompt: DefaultPrompt,
		Value:  drawable,
	}
}

func (i *InputLine) ToDrawable() drawable.Drawable {
	return toDrawable(i)
}

type inputLineDrawable struct {
	initialized bool
	status      bool
	prompt      string
	input       drawable.Drawable
}

func toDrawable(input *InputLine) drawable.Drawable {
	drw := inputLineDrawable{
		initialized: false,
		status:      true,
		prompt:      input.Prompt,
		input:       input.Value,
	}
	return drawable.Drawable{
		Init: drw.init,
		Draw: drw.draw,
	}
}

func (d *inputLineDrawable) init(size terminal.Winsize) {
	d.initialized = true

	d.input.Init(size)
}

func (d *inputLineDrawable) draw() ([]text.Line, bool) {
	assert.True(d.initialized, "the drawable should be initialized before draw")

	lines := make([]text.Line, 0)

	if !d.status {
		return lines, false
	}

	content := true
	for content {
		result, status := d.input.Draw()
		content = status
		lines = append(lines, result...)
	}

	if len(lines) == 0 {
		line := text.LineFromString(d.prompt)
		return []text.Line{line}, false
	}

	prompt := text.FragmentsFromString(d.prompt + marker.DefaultPaddingText)
	lines[0].Text = append(prompt, lines[0].Text...)

	return append([]text.Line{text.LineFromString("")}, lines...), false
}
