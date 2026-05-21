package inline

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

const Name = "inline_unit"

type InlineUnit struct {
	loaded     bool
	separator  string
	units      []drawable.Unit
}

func New(units ...drawable.Unit) *InlineUnit {
	return &InlineUnit{
		loaded:     false,
		separator:  "",
		units:      units,
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
}

func (d *InlineUnit) wipe() {}

func (d *InlineUnit) draw(size winsize.Winsize) ([]text.Line, bool) {
	assert.True(d.loaded, drawable.MessageInitialized)

	lines := d.drawChildren(size)

	return d.joinChildren(lines), false
}

func (d *InlineUnit) drawChildren(size winsize.Winsize) []text.Line {
	lines := make([]text.Line, 0)

	if len(d.units) == 0 {
		return lines
	}

	index := 0

	focus := d.units[index]
	focus.Drawable.Init()

	for {
		result, status := focus.Drawable.Draw(size)
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
