package inputline

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/utils/drain"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/marker"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

const Name = "input_line_unit"

type InputLineUnit struct {
	loaded bool
	status bool
	prompt string
	unit   drawable.Unit
}

func New(unit drawable.Unit) *InputLineUnit {
	return &InputLineUnit{
		loaded: false,
		status: true,
		prompt: marker.DefaultInputLinePrompt,
		unit:   unit,
	}
}

func UnitFromUnit(unit drawable.Unit) drawable.Unit {
	return New(unit).ToUnit()
}

func (d *InputLineUnit) ToUnit() drawable.Unit {
	return drawable.NewBuilder().
		Name(Name).
		Init(d.init).
		Wipe(d.unit.Drawable.Wipe).
		Draw(d.draw).
		ToUnit()
}

func (d *InputLineUnit) init() {
	d.loaded = true

	d.unit.Drawable.Init()
}

func (d *InputLineUnit) draw(size winsize.Winsize) ([]text.Line, bool) {
	assert.True(d.loaded, drawable.MessageInitialized)

	if size.Rows == 0 {
		return make([]text.Line, 0), false
	}

	lines, _ := drain.UnitLazy(size, d.unit)
	if len(lines) == 0 {
		line := text.NewLine(d.prompt)
		return []text.Line{*line}, false
	}

	prompt := text.FragmentsFromString(d.prompt + marker.DefaultPaddingText)
	lines[0].Text = append(prompt, lines[0].Text...)

	return lines, false
}
