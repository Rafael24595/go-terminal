package checkmenu

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/primitive/line"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/spatial/justify"
	"github.com/Rafael24595/go-reacterm-core/engine/model/input"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/marker"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

const Name = "check_menu_unit"

type CheckMenuUnit struct {
	initialized  bool
	meta         marker.CheckMeta
	distribution style.Distribution
	writeMode    bool
	options      []input.CheckOption
	cursor       uint16
	unit         drawable.Unit
}

func New(options []input.CheckOption) *CheckMenuUnit {
	clone := make([]input.CheckOption, len(options))
	copy(clone, options)

	return &CheckMenuUnit{
		initialized:  false,
		meta:         marker.BracketsCheck,
		distribution: style.DefaultDistribution,
		options:      clone,
		cursor:       0,
		unit:         drawable.Unit{},
	}
}

func UnitFromOptions(options []input.CheckOption) drawable.Unit {
	return New(options).ToUnit()
}

func (d *CheckMenuUnit) Meta(meta marker.CheckMeta) *CheckMenuUnit {
	d.meta = meta
	return d
}

func (d *CheckMenuUnit) Distribution(distribution style.Distribution) *CheckMenuUnit {
	d.distribution = distribution
	return d
}

func (d *CheckMenuUnit) WriteMode(writeMode bool) *CheckMenuUnit {
	d.writeMode = writeMode
	return d
}

func (d *CheckMenuUnit) Cursor(cursor uint16) *CheckMenuUnit {
	d.cursor = cursor
	return d
}

func (d *CheckMenuUnit) ToUnit() drawable.Unit {
	return drawable.NewBuilder().
		Name(Name).
		Init(d.init).
		Wipe(d.wipe).
		Draw(d.draw).
		ToUnit()
}

func (d *CheckMenuUnit) init() {
	d.initialized = true

	opts := d.addStyles()

	switch d.distribution.Direction {
	case style.Vertical:
		d.unit = d.makeVertical(opts)
	case style.Horizontal:
		d.unit = d.makeHorizontal(opts)
	default:
		assert.Unreachable("undefined direction %d", d.distribution.Direction)

		d.unit = d.makeVertical(opts)
	}

	d.unit.Drawable.Init()
}

func (d *CheckMenuUnit) wipe() {
	if d.unit.Drawable.Wipe == nil {
		return
	}
	d.unit.Drawable.Wipe()
}

func (d *CheckMenuUnit) makeVertical(opts []text.Fragment) drawable.Unit {
	lines := make([]text.Line, len(opts))
	for i := range opts {
		lines[i] = *text.LineFromFragments(opts[i])
	}
	return line.UnitFromLines(lines...)
}

func (d *CheckMenuUnit) makeHorizontal(opts []text.Fragment) drawable.Unit {
	return justify.New(opts).
		Justify(d.distribution.Justify).
		MaxOpts(d.distribution.Limit).
		ToUnit()
}

func (d *CheckMenuUnit) addStyles() []text.Fragment {
	frags := make([]text.Fragment, len(d.options))

	for i := range frags {
		status := d.meta.Unchecked
		if d.options[i].Status {
			status = d.meta.Checked
		}

		label := d.options[i].Label.Text
		if len(label) > 0 {
			label = marker.DefaultPaddingText + label
		}

		frags[i] = *text.NewFragment(d.meta.Open + status + d.meta.Close + label).
			CopyMeta(&d.options[i].Label)

		if d.writeMode && i == int(d.cursor) {
			frags[i].AddAtom(style.AtmSelect, style.AtmFocus)
		}
	}

	return frags
}

func (d *CheckMenuUnit) draw(size winsize.Winsize) ([]text.Line, bool) {
	assert.True(d.initialized, drawable.MessageInitialized)

	return d.unit.Drawable.Draw(size)
}
