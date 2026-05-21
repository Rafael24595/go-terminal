package pipeline

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

const Name = "pipeline_unit"

type InitTransformer func(winsize.Winsize, drawable.Unit) drawable.Unit
type DrawTransformer func(winsize.Winsize, drawable.Unit) ([]text.Line, bool)
type DataTransformer func(winsize.Winsize, drawable.Unit, []text.Line, bool) ([]text.Line, bool)

type PipelineUnit struct {
	loaded    bool
	unit      drawable.Unit
	initSteps []InitTransformer
	drawStep  DrawTransformer
	dataSteps []DataTransformer
}

func New(unit drawable.Unit) *PipelineUnit {
	return &PipelineUnit{
		loaded:    false,
		unit:      unit,
		initSteps: make([]InitTransformer, 0),
		drawStep:  nil,
		dataSteps: make([]DataTransformer, 0),
	}
}

func DrawToUnit(unit drawable.Unit, step DrawTransformer) drawable.Unit {
	return New(unit).
		SetDrawStep(step).
		ToUnit()
}

func (d *PipelineUnit) PushInitSteps(steps ...InitTransformer) *PipelineUnit {
	if d.loaded {
		assert.Unreachable(drawable.MessageNewElement)
		return d
	}

	d.initSteps = append(d.initSteps, steps...)
	return d
}

func (d *PipelineUnit) SetDrawStep(step DrawTransformer) *PipelineUnit {
	if d.loaded {
		assert.Unreachable(drawable.MessageNewElement)
		return d
	}

	d.drawStep = step
	return d
}

func (d *PipelineUnit) PushDataSteps(steps ...DataTransformer) *PipelineUnit {
	if d.loaded {
		assert.Unreachable(drawable.MessageNewElement)
		return d
	}

	d.dataSteps = append(d.dataSteps, steps...)
	return d
}

func (d *PipelineUnit) ToUnit() drawable.Unit {
	if d.isAnemic() {
		return d.unit
	}

	return drawable.NewBuilder().
		Name(Name).
		MergeTags(d.unit.Tags).
		Init(d.init).
		Wipe(d.wipe).
		Draw(d.draw).
		ToUnit()
}

func (d *PipelineUnit) isAnemic() bool {
	if d.drawStep != nil {
		return false
	}

	if len(d.initSteps) > 0 {
		return false
	}

	return len(d.dataSteps) == 0
}

func (d *PipelineUnit) init() {
	d.loaded = true

	d.unit.Drawable.Init()
}

func (d *PipelineUnit) wipe() {
	if d.unit.Drawable.Wipe == nil {
		return
	}
	d.unit.Drawable.Wipe()
}

func (d *PipelineUnit) draw(size winsize.Winsize) ([]text.Line, bool) {
	assert.True(d.loaded, drawable.MessageInitialized)

	for _, s := range d.initSteps {
		d.unit = s(size, d.unit)
	}

	draw := d.unit.Drawable.Draw
	if d.drawStep != nil {
		draw = func(size winsize.Winsize) ([]text.Line, bool) {
			return d.drawStep(size, d.unit)
		}
	}

	lines, status := draw(size)
	for _, s := range d.dataSteps {
		lines, status = s(size, d.unit, lines, status)
	}

	return lines, status
}
