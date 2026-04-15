package mapdrw

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-terminal/engine/layout/drawable"
	"github.com/Rafael24595/go-terminal/engine/render/text"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

const NameMapDrawable = "MapDrawable"

type drawInputPred func(terminal.Winsize) terminal.Winsize
type drawOutputPred func(terminal.Winsize, drawable.Drawable) ([]text.Line, bool)

type MapDrawable struct {
	loaded        bool
	drawInputMap  drawInputPred
	drawOutputMap drawOutputPred
	drawable      drawable.Drawable
}

func NewMapDrawable(drw drawable.Drawable) *MapDrawable {
	return &MapDrawable{
		loaded:        false,
		drawInputMap:  nil,
		drawOutputMap: nil,
		drawable:      drw,
	}
}

func (d *MapDrawable) SetDrawInputMap(pred drawInputPred) *MapDrawable {
	if d.loaded {
		assert.Unreachable(drawable.MessageNewElement)
		return d
	}

	d.drawInputMap = pred
	return d
}

func (d *MapDrawable) SetDrawOutputMap(pred drawOutputPred) *MapDrawable {
	if d.loaded {
		assert.Unreachable(drawable.MessageNewElement)
		return d
	}

	d.drawOutputMap = pred
	return d
}

func (d *MapDrawable) ToDrawable() drawable.Drawable {
	nilDrawMap := d.drawInputMap == nil && d.drawOutputMap == nil
	if nilDrawMap {
		return d.drawable
	}

	return drawable.Drawable{
		Name: NameMapDrawable,
		Code: d.drawable.Code,
		Tags: d.drawable.Tags,
		Init: d.init,
		Wipe: d.wipe,
		Draw: d.draw,
	}
}

func (d *MapDrawable) init() {
	d.loaded = true

	d.drawable.Init()
}

func (d *MapDrawable) wipe() {
	if d.drawable.Wipe == nil {
		return
	}
	d.drawable.Wipe()
}

func (d *MapDrawable) draw(size terminal.Winsize) ([]text.Line, bool) {
	assert.True(d.loaded, drawable.MessageInitialized)

	mapSize := size
	if d.drawInputMap != nil {
		mapSize = d.drawInputMap(mapSize)
	}

	if d.drawOutputMap == nil {
		return d.drawable.Draw(mapSize)
	}

	return d.drawOutputMap(mapSize, d.drawable)
}
