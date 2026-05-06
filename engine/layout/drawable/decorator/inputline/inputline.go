package inputline

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/set"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/marker"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

const Name = "input_line_drawable"

type InputLineDrawable struct {
	loaded   bool
	status   bool
	prompt   string
	drawable drawable.Drawable
}

func New(drawable drawable.Drawable) *InputLineDrawable {
	return &InputLineDrawable{
		loaded:   false,
		status:   true,
		prompt:   marker.DefaultInputLinePrompt,
		drawable: drawable,
	}
}

func DrawableFromDrawable(drawable drawable.Drawable) drawable.Drawable {
	return New(drawable).ToDrawable()
}

func (d *InputLineDrawable) ToDrawable() drawable.Drawable {
	return drawable.Drawable{
		Name: Name,
		Code: "",
		Tags: make(set.Set[string]),
		Init: d.init,
		Wipe: d.drawable.Wipe,
		Draw: d.draw,
	}
}

func (d *InputLineDrawable) init() {
	d.loaded = true

	d.drawable.Init()
}

func (d *InputLineDrawable) draw(size winsize.Winsize) ([]text.Line, bool) {
	assert.True(d.loaded, drawable.MessageInitialized)

	if size.Rows == 0 {
		return make([]text.Line, 0), false
	}

	lines, _ := drawable.DrainDrawable(size, d.drawable, true)
	if len(lines) == 0 {
		line := text.NewLine(d.prompt)
		return []text.Line{*line}, false
	}

	prompt := text.FragmentsFromString(d.prompt + marker.DefaultPaddingText)
	lines[0].Text = append(prompt, lines[0].Text...)

	return lines, false
}
