package indexmenu

import (
	"strconv"

	"github.com/Rafael24595/go-terminal/engine/core/assert"
	"github.com/Rafael24595/go-terminal/engine/core/drawable"
	"github.com/Rafael24595/go-terminal/engine/core/drawable/line"
	"github.com/Rafael24595/go-terminal/engine/core/marker"
	"github.com/Rafael24595/go-terminal/engine/core/style"
	"github.com/Rafael24595/go-terminal/engine/core/text"
	"github.com/Rafael24595/go-terminal/engine/helper"
	"github.com/Rafael24595/go-terminal/engine/helper/math"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

type IndexMenuDrawable struct {
	initialized bool
	index       marker.IndexMeta
	options     []text.Fragment
	cursor      uint
	drawable    drawable.Drawable
}

func NewIndexMenuDrawable(index marker.IndexMeta, options []text.Fragment) *IndexMenuDrawable {
	clone := make([]text.Fragment, len(options))
	copy(clone, options)

	return &IndexMenuDrawable{
		initialized: false,
		index:       index,
		options:     clone,
		cursor:      0,
		drawable:    drawable.Drawable{},
	}
}

func TextIndexMenuFromData(index marker.IndexMeta, options []text.Fragment) drawable.Drawable {
	return NewIndexMenuDrawable(index, options).ToDrawable()
}

func (d *IndexMenuDrawable) Cursor(cursor uint) *IndexMenuDrawable {
	d.cursor = cursor
	return d
}

func (d *IndexMenuDrawable) ToDrawable() drawable.Drawable {
	return drawable.Drawable{
		Init: d.init,
		Draw: d.draw,
	}
}

func (d *IndexMenuDrawable) init(size terminal.Winsize) {
	d.initialized = true

	lines := make([]text.Line, 0)

	digits := math.Digits(len(d.options))

	for i, o := range d.options {
		styl := style.AtmNone
		if i == int(d.cursor) {
			styl = style.AtmFocus
		}

		styledLine := text.LineFromFragments(
			text.EmptyFragment().
				AddSpec(style.SpecPaddingLeft(2)),
			d.makeIndex(i, int(digits)),
			text.NewFragment(" "),
			text.NewFragment(o.Text).
				AddAtom(styl),
		)

		lines = append(lines, styledLine)
	}

	drawable := line.EagerDrawableFromLines(lines...)
	drawable.Init(size)

	d.drawable = drawable
}

func (d *IndexMenuDrawable) makeIndex(cursor, digits int) text.Fragment {
	if d.index.Kind == marker.Numeric {
		txt := helper.Right(strconv.Itoa(cursor+1), digits)
		index := text.NewFragment(txt + ".- ")
		if cursor == int(d.cursor) {
			index.Atom |= style.AtmBold
		}
		return index
	}

	if d.index.Kind == marker.Alphabetic {
		txt := helper.Right(helper.NumberToAlpha(cursor), digits)
		index := text.NewFragment(txt + ".- ")
		if cursor == int(d.cursor) {
			index.Atom |= style.AtmBold
		}
		return index
	}

	index := d.index.Index
	if cursor == int(d.cursor) {
		index = d.index.Cursor
	}

	return text.NewFragment(index)
}

func (d *IndexMenuDrawable) draw() ([]text.Line, bool) {
	assert.True(d.initialized, "the drawable should be initialized before draw")

	return d.drawable.Draw()
}
