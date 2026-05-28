package inputline

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/transform/drain"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/marker"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"

	drawable_drain "github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/drain"
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
		prompt: marker.DefaultPromptText,
		unit:   unit,
	}
}

func Wrap(unit drawable.Unit) drawable.Unit {
	return New(unit).ToUnit()
}

func FromFragment(frag text.Fragment) drawable.Unit {
	return Wrap(drawable_drain.UnitFromFragments(frag))
}

func (u *InputLineUnit) ToUnit() drawable.Unit {
	return drawable.NewBuilder().
		Name(Name).
		Init(u.init).
		Wipe(u.unit.Drawable.Wipe).
		Draw(u.draw).
		ToUnit()
}

func (u *InputLineUnit) init() {
	u.loaded = true

	u.unit.Drawable.Init()
}

func (u *InputLineUnit) draw(size winsize.Winsize) ([]text.Line, bool) {
	assert.True(u.loaded, drawable.MessageInitialized)

	if size.Rows == 0 {
		return make([]text.Line, 0), false
	}

	lines, _ := drain.UnitLazy(size, u.unit)
	if len(lines) == 0 {
		line := text.NewLine(u.prompt)
		return []text.Line{*line}, false
	}

	prompt := text.FragmentsFromString(u.prompt + marker.DefaultPaddingText)
	lines[0].Text = append(prompt, lines[0].Text...)

	return lines, false
}
