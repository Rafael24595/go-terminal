package indexmenu

import (
	"strconv"

	assert "github.com/Rafael24595/go-assert/assert/runtime"
	
	"github.com/Rafael24595/go-terminal/engine/helper"
	"github.com/Rafael24595/go-terminal/engine/helper/math"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable/line"
	"github.com/Rafael24595/go-terminal/engine/render/marker"
	"github.com/Rafael24595/go-terminal/engine/render/style"
	"github.com/Rafael24595/go-terminal/engine/render/text"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

const NameIndexMenuDrawable = "IndexMenuDrawable"

type IndexMenuDrawable struct {
	initialized bool
	meta        marker.IndexMeta
	options     []text.Fragment
	cursor      uint
	drawable    drawable.Drawable
}

func NewIndexMenuDrawable(options []text.Fragment) *IndexMenuDrawable {
	clone := make([]text.Fragment, len(options))
	copy(clone, options)

	return &IndexMenuDrawable{
		initialized: false,
		meta:        marker.HyphenIndex,
		options:     clone,
		cursor:      0,
		drawable:    drawable.Drawable{},
	}
}

func TextIndexMenuFromData(options []text.Fragment) drawable.Drawable {
	return NewIndexMenuDrawable(options).ToDrawable()
}

func (d *IndexMenuDrawable) Meta(meta marker.IndexMeta) *IndexMenuDrawable {
	d.meta = meta
	return d
}

func (d *IndexMenuDrawable) Cursor(cursor uint) *IndexMenuDrawable {
	d.cursor = cursor
	return d
}

func (d *IndexMenuDrawable) ToDrawable() drawable.Drawable {
	return drawable.Drawable{
		Name: NameIndexMenuDrawable,
		Init: d.init,
		Draw: d.draw,
	}
}

func (d *IndexMenuDrawable) init(size terminal.Winsize) {
	d.initialized = true

	lines := make([]text.Line, 0)

	digits := math.Digits(len(d.options))

	for i, o := range d.options {
		focs := style.AtmNone
		if i == int(d.cursor) {
			focs = style.AtmFocus
		}

		padd := text.EmptyFragment().
			AddSpec(style.SpecPaddingLeft(2))
		indx := d.makeIndex(i, int(digits))
		spac := text.NewFragment(marker.DefaultPaddingText)
		mark := text.NewFragment(o.Text).
			AddAtom(focs)

		lines = append(lines,
			text.LineFromFragments(padd, indx, spac, mark),
		)
	}

	drawable := line.EagerDrawableFromLines(lines...)
	drawable.Init(size)

	d.drawable = drawable
}

func (d *IndexMenuDrawable) makeIndex(cursor, digits int) text.Fragment {
	if d.meta.Kind == marker.Numeric {
		txt := helper.Right(strconv.Itoa(cursor+1), digits)
		index := text.NewFragment(txt + ".- ")
		if cursor == int(d.cursor) {
			index.Atom |= style.AtmBold
		}
		return index
	}

	if d.meta.Kind == marker.Alphabetic {
		txt := helper.Right(helper.NumberToAlpha(cursor), digits)
		index := text.NewFragment(txt + ".- ")
		if cursor == int(d.cursor) {
			index.Atom |= style.AtmBold
		}
		return index
	}

	index := d.meta.Index
	if cursor == int(d.cursor) {
		index = d.meta.Cursor
	}

	return text.NewFragment(index)
}

func (d *IndexMenuDrawable) draw() ([]text.Line, bool) {
	assert.True(d.initialized, "the drawable should be initialized before draw")

	return d.drawable.Draw()
}
