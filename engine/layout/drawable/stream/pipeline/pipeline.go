package pipeline

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

const Name = "pipeline_drawable"

type InitTransformer func(winsize.Winsize, drawable.Drawable) drawable.Drawable
type DrawTransformer func(winsize.Winsize, drawable.Drawable) ([]text.Line, bool)
type DataTransformer func(winsize.Winsize, drawable.Drawable, []text.Line, bool) ([]text.Line, bool)

type PipelineDrawable struct {
	loaded    bool
	drawable  drawable.Drawable
	initSteps []InitTransformer
	drawStep  DrawTransformer
	dataSteps []DataTransformer
}

func New(drawable drawable.Drawable) *PipelineDrawable {
	return &PipelineDrawable{
		loaded:    false,
		drawable:  drawable,
		initSteps: make([]InitTransformer, 0),
		drawStep:  nil,
		dataSteps: make([]DataTransformer, 0),
	}
}

func DrawToDrawable(drawable drawable.Drawable, step DrawTransformer) drawable.Drawable {
	return New(drawable).
		SetDrawStep(step).
		ToDrawable()
}

func (d *PipelineDrawable) PushInitSteps(steps ...InitTransformer) *PipelineDrawable {
	if d.loaded {
		assert.Unreachable(drawable.MessageNewElement)
		return d
	}

	d.initSteps = append(d.initSteps, steps...)
	return d
}

func (d *PipelineDrawable) SetDrawStep(step DrawTransformer) *PipelineDrawable {
	if d.loaded {
		assert.Unreachable(drawable.MessageNewElement)
		return d
	}

	d.drawStep = step
	return d
}

func (d *PipelineDrawable) PushDataSteps(steps ...DataTransformer) *PipelineDrawable {
	if d.loaded {
		assert.Unreachable(drawable.MessageNewElement)
		return d
	}

	d.dataSteps = append(d.dataSteps, steps...)
	return d
}

func (d *PipelineDrawable) ToDrawable() drawable.Drawable {
	if len(d.initSteps) == 0 &&
		d.drawStep == nil &&
		len(d.dataSteps) == 0 {
		return d.drawable
	}

	return drawable.Drawable{
		Name: Name,
		Code: d.drawable.Code,
		Tags: d.drawable.Tags,
		Init: d.init,
		Wipe: d.wipe,
		Draw: d.draw,
	}
}

func (d *PipelineDrawable) init() {
	d.loaded = true

	d.drawable.Init()
}

func (d *PipelineDrawable) wipe() {
	if d.drawable.Wipe == nil {
		return
	}
	d.drawable.Wipe()
}

func (d *PipelineDrawable) draw(size winsize.Winsize) ([]text.Line, bool) {
	assert.True(d.loaded, drawable.MessageInitialized)

	for _, s := range d.initSteps {
		d.drawable = s(size, d.drawable)
	}

	draw := d.drawable.Draw
	if d.drawStep != nil {
		draw = func(size winsize.Winsize) ([]text.Line, bool) {
			return d.drawStep(size, d.drawable)
		}
	}

	lines, status := draw(size)
	for _, s := range d.dataSteps {
		lines, status = s(size, d.drawable, lines, status)
	}

	return lines, status
}
