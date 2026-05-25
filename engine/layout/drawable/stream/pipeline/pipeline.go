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

func (u *PipelineUnit) PushInitSteps(steps ...InitTransformer) *PipelineUnit {
	if u.loaded {
		assert.Unreachable(drawable.MessageNewElement)
		return u
	}

	u.initSteps = append(u.initSteps, steps...)
	return u
}

func (u *PipelineUnit) SetDrawStep(step DrawTransformer) *PipelineUnit {
	if u.loaded {
		assert.Unreachable(drawable.MessageNewElement)
		return u
	}

	u.drawStep = step
	return u
}

func (u *PipelineUnit) PushDataSteps(steps ...DataTransformer) *PipelineUnit {
	if u.loaded {
		assert.Unreachable(drawable.MessageNewElement)
		return u
	}

	u.dataSteps = append(u.dataSteps, steps...)
	return u
}

func (u *PipelineUnit) ToUnit() drawable.Unit {
	if u.isAnemic() {
		return u.unit
	}

	return drawable.NewBuilder().
		Name(Name).
		MergeTags(u.unit.Tags).
		Init(u.init).
		Wipe(u.wipe).
		Draw(u.draw).
		ToUnit()
}

func (u *PipelineUnit) isAnemic() bool {
	if u.drawStep != nil {
		return false
	}

	if len(u.initSteps) > 0 {
		return false
	}

	return len(u.dataSteps) == 0
}

func (u *PipelineUnit) init() {
	u.loaded = true

	u.unit.Drawable.Init()
}

func (u *PipelineUnit) wipe() {
	if u.unit.Drawable.Wipe == nil {
		return
	}
	u.unit.Drawable.Wipe()
}

func (u *PipelineUnit) draw(size winsize.Winsize) ([]text.Line, bool) {
	assert.True(u.loaded, drawable.MessageInitialized)

	for _, s := range u.initSteps {
		u.unit = s(size, u.unit)
	}

	draw := u.unit.Drawable.Draw
	if u.drawStep != nil {
		draw = func(size winsize.Winsize) ([]text.Line, bool) {
			return u.drawStep(size, u.unit)
		}
	}

	lines, status := draw(size)
	for _, s := range u.dataSteps {
		lines, status = s(size, u.unit, lines, status)
	}

	return lines, status
}
