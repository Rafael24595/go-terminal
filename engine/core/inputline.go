package core

import (
	"github.com/Rafael24595/go-terminal/engine/core/assert"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

const DefaultPrompt = ">"

type InputLine struct {
	Prompt string
	Value  Drawable
}

func NewInputLine(drawable Drawable) *InputLine {
	return &InputLine{
		Prompt: DefaultPrompt,
		Value:  drawable,
	}
}

func (i *InputLine) ToDrawable() Drawable {
	return toDrawable(i)
}

type inputLineDrawable struct {
	initialized bool
	status      bool
	prompt      string
	input       Drawable
}

func toDrawable(input *InputLine) Drawable {
	drawable := inputLineDrawable{
		initialized: false,
		status:      true,
		prompt:      input.Prompt,
		input:       input.Value,
	}
	return Drawable{
		Init: drawable.init,
		Draw: drawable.draw,
	}
}

func (d *inputLineDrawable) init(size terminal.Winsize) {
	d.initialized = true

	d.input.Init(size)
}

func (d *inputLineDrawable) draw() ([]Line, bool) {
	assert.True(d.initialized, "the drawable should be initialized before draw")

	lines := make([]Line, 0)

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
		line := LineFromString(d.prompt)
		return []Line{line}, false
	}

	prompt := FragmentsFromString(d.prompt + " ")
	lines[0].Text = append(prompt, lines[0].Text...)

	return append([]Line{ LineFromString("") }, lines...), false
}
