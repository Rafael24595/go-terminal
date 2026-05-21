package inline

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/drain"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

const Name = "inline_unit"

//TODO: Remove lazy init.
type InlineUnit struct {
	loaded     bool
	lazyLoaded bool
	size       winsize.Winsize
	separator  string
	units      []drawable.Unit
	unit       drawable.Unit
}

func New(units ...drawable.Unit) *InlineUnit {
	return &InlineUnit{
		loaded:     false,
		lazyLoaded: false,
		size:       winsize.Winsize{},
		separator:  "",
		units:      units,
		unit:       drawable.Unit{},
	}
}

func UnitFromUnits(units ...drawable.Unit) drawable.Unit {
	return New(units...).ToUnit()
}

func (d *InlineUnit) Separator(separator string) *InlineUnit {
	d.separator = separator
	return d
}

func (d *InlineUnit) ToUnit() drawable.Unit {
	return drawable.NewBuilder().
		Name(Name).
		Init(d.init).
		Wipe(d.wipe).
		Draw(d.draw).
		ToUnit()
}

func (d *InlineUnit) init() {
	d.loaded = true
	d.lazyLoaded = false
}

func (d *InlineUnit) wipe() {
	d.lazyLoaded = false
}

func (d *InlineUnit) lazyInit(size winsize.Winsize) {
	if d.lazyLoaded {
		return
	}

	d.lazyLoaded = true

	d.size = size

	lines := d.drawChildren()
	join := d.joinChildren(lines)

	d.unit = drain.UnitFromLines(join...)

	d.unit.Drawable.Init()
}

func (d *InlineUnit) draw(size winsize.Winsize) ([]text.Line, bool) {
	assert.True(d.loaded, drawable.MessageInitialized)

	d.lazyInit(size)

	return d.unit.Drawable.Draw(size)
}

func (d *InlineUnit) drawChildren() []text.Line {
	lines := make([]text.Line, 0)

	if len(d.units) == 0 {
		return lines
	}

	index := 0

	focus := d.units[index]
	focus.Drawable.Init()

	for {
		result, status := focus.Drawable.Draw(d.size)
		if len(result) > 0 {
			lines = append(lines, result...)
		}

		if status && len(result) == 0 {
			continue
		}

		index += 1
		if index >= len(d.units) {
			break
		}

		focus = d.units[index]
		focus.Drawable.Init()
	}

	return lines
}

func (d *InlineUnit) joinChildren(lines []text.Line) []text.Line {
	if len(lines) == 0 {
		return []text.Line{}
	}

	merged := text.EmptyLine()

	var separator *text.Fragment
	if d.separator != "" {
		frag := text.NewFragment(d.separator)
		separator = frag
	}

	for i, line := range lines {
		frags := line.Text
		if d.separator != "" && i < len(lines)-1 {
			frags = append(frags, *separator)
		}

		merged.PushFragments(frags...).
			CopyMeta(&line)
	}

	return []text.Line{*merged}
}
