package checkmenu

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-terminal/engine/layout/drawable"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable/spatial/justify"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable/primitive/line"
	"github.com/Rafael24595/go-terminal/engine/model/input"
	"github.com/Rafael24595/go-terminal/engine/render/marker"
	"github.com/Rafael24595/go-terminal/engine/render/style"
	"github.com/Rafael24595/go-terminal/engine/render/text"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

const NameCheckMenuDrawable = "CheckMenuDrawable"

type CheckMenuDrawable struct {
	initialized  bool
	meta         marker.CheckMeta
	distribution style.Distribution
	writeMode    bool
	options      []input.CheckOption
	cursor       uint
	drawable     drawable.Drawable
}

func NewCheckMenuDrawable(options []input.CheckOption) *CheckMenuDrawable {
	clone := make([]input.CheckOption, len(options))
	copy(clone, options)

	return &CheckMenuDrawable{
		initialized:  false,
		meta:         marker.BracketsCheck,
		distribution: style.DefaultDistribution,
		options:      clone,
		cursor:       0,
		drawable:     drawable.Drawable{},
	}
}

func CheckMenuDrawableOptions(options []input.CheckOption) drawable.Drawable {
	return NewCheckMenuDrawable(options).ToDrawable()
}

func (d *CheckMenuDrawable) Meta(meta marker.CheckMeta) *CheckMenuDrawable {
	d.meta = meta
	return d
}

func (d *CheckMenuDrawable) Distribution(distribution style.Distribution) *CheckMenuDrawable {
	d.distribution = distribution
	return d
}

func (d *CheckMenuDrawable) WriteMode(writeMode bool) *CheckMenuDrawable {
	d.writeMode = writeMode
	return d
}

func (d *CheckMenuDrawable) Cursor(cursor uint) *CheckMenuDrawable {
	d.cursor = cursor
	return d
}

func (d *CheckMenuDrawable) ToDrawable() drawable.Drawable {
	return drawable.Drawable{
		Name: NameCheckMenuDrawable,
		Code: d.drawable.Code,
		Tags: d.drawable.Tags,
		Init: d.init,
		Wipe: d.wipe,
		Draw: d.draw,
	}
}

func (d *CheckMenuDrawable) init() {
	d.initialized = true

	opts := d.addStyles()

	switch d.distribution.Direction {
	case style.Vertical:
		d.drawable = d.makeVertical(opts)
	case style.Horizontal:
		d.drawable = d.makeHorizontal(opts)
	default:
		assert.Unreachable("undefined direction %d", d.distribution.Direction)

		d.drawable = d.makeVertical(opts)
	}

	d.drawable.Init()
}

func (d *CheckMenuDrawable) wipe() {
	if d.drawable.Wipe == nil {
		return
	}
	d.drawable.Wipe()
}

func (d *CheckMenuDrawable) makeVertical(opts []text.Fragment) drawable.Drawable {
	lines := make([]text.Line, len(opts))
	for i := range opts {
		lines[i] = *text.LineFromFragments(opts[i])
	}
	return line.LineDrawableFromLines(lines...)
}

func (d *CheckMenuDrawable) makeHorizontal(opts []text.Fragment) drawable.Drawable {
	return justify.NewJustifyDrawable(opts).
		Justify(d.distribution.Justify).
		MaxOpts(d.distribution.Limit).
		ToDrawable()
}

func (d *CheckMenuDrawable) addStyles() []text.Fragment {
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

func (d *CheckMenuDrawable) draw(size terminal.Winsize) ([]text.Line, bool) {
	assert.True(d.initialized, drawable.MessageInitialized)

	return d.drawable.Draw(size)
}
