package check

import (
	"github.com/Rafael24595/go-terminal/engine/core/assert"
	"github.com/Rafael24595/go-terminal/engine/core/drawable"
	"github.com/Rafael24595/go-terminal/engine/core/drawable/justify"
	"github.com/Rafael24595/go-terminal/engine/core/drawable/line"
	"github.com/Rafael24595/go-terminal/engine/core/input"
	"github.com/Rafael24595/go-terminal/engine/core/marker"
	"github.com/Rafael24595/go-terminal/engine/core/style"
	"github.com/Rafael24595/go-terminal/engine/core/text"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

type ChecMenukDrawable struct {
	initialized  bool
	meta         marker.CheckMeta
	distribution style.Distribution
	writeMode    bool
	options      []input.CheckOption
	cursor       uint
	drawable     drawable.Drawable
}

func NewCheckMenuDrawable(options []input.CheckOption) *ChecMenukDrawable {
	clone := make([]input.CheckOption, len(options))
	copy(clone, options)

	return &ChecMenukDrawable{
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

func (d *ChecMenukDrawable) Meta(meta marker.CheckMeta) *ChecMenukDrawable {
	d.meta = meta
	return d
}

func (d *ChecMenukDrawable) Distribution(distribution style.Distribution) *ChecMenukDrawable {
	d.distribution = distribution
	return d
}

func (d *ChecMenukDrawable) Cursor(cursor uint) *ChecMenukDrawable {
	d.cursor = cursor
	return d
}

func (d *ChecMenukDrawable) ToDrawable() drawable.Drawable {
	return drawable.Drawable{
		Init: d.init,
		Draw: d.draw,
	}
}

func (d *ChecMenukDrawable) init(size terminal.Winsize) {
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
}

func (d *ChecMenukDrawable) makeVertical(opts []text.Fragment) drawable.Drawable {
	lines := make([]text.Line, len(opts))
	for i := range opts {
		lines[i] = text.LineFromFragments(opts[i])
	}
	return line.LazyDrawableFromLines(lines...)
}

func (d *ChecMenukDrawable) makeHorizontal(opts []text.Fragment) drawable.Drawable {
	return justify.NewJustifyDrawable(opts).
		Justify(d.distribution.Justify).
		Limit(d.distribution.Limit).
		ToDrawable()
}

func (d *ChecMenukDrawable) addStyles() []text.Fragment {
	frags := make([]text.Fragment, len(d.options))

	for i := range frags {
		status := d.meta.Unchecked
		if d.options[i].Status {
			status = d.meta.Checked
		}

		label := d.options[i].Label.Text
		if len(label) > 0 {
			label = " " + label
		}

		frags[i] = text.EmptyFragmentFrom(d.options[i].Label)
		frags[i].Text = d.meta.Close + status + d.meta.Close + label

		if i == int(d.cursor) {
			frags[i] = frags[i].
				AddAtom(style.AtmSelect, style.AtmFocus)
		}
	}

	return frags
}

func (d *ChecMenukDrawable) draw() ([]text.Line, bool) {
	assert.True(d.initialized, "the drawable should be initialized before draw")

	return d.drawable.Draw()
}
