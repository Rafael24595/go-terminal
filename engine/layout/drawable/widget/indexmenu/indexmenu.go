package indexmenu

import (
	"strconv"

	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/helper"
	"github.com/Rafael24595/go-reacterm-core/engine/helper/math"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/drain"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/marker"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

const Name = "index_menu_unit"

type IndexMenuUnit struct {
	loaded  bool
	meta    marker.IndexMeta
	options []text.Fragment
	cursor  uint16
	unit    drawable.Unit
}

func New(options []text.Fragment) *IndexMenuUnit {
	clone := make([]text.Fragment, len(options))
	copy(clone, options)

	return &IndexMenuUnit{
		loaded:  false,
		meta:    marker.HyphenIndex,
		options: clone,
		cursor:  0,
		unit:    drawable.Unit{},
	}
}

func UnitFromOptions(options []text.Fragment) drawable.Unit {
	return New(options).ToUnit()
}

func (d *IndexMenuUnit) Meta(meta marker.IndexMeta) *IndexMenuUnit {
	d.meta = meta
	return d
}

func (d *IndexMenuUnit) Cursor(cursor uint16) *IndexMenuUnit {
	d.cursor = cursor
	return d
}

func (d *IndexMenuUnit) ToUnit() drawable.Unit {
	return drawable.NewBuilder().
		Name(Name).
		Init(d.init).
		Wipe(d.wipe).
		Draw(d.draw).
		ToUnit()
}

func (d *IndexMenuUnit) init() {
	d.loaded = true

	lines := make([]text.Line, 0)

	digits := math.Digits(len(d.options))

	for i, o := range d.options {
		focs := style.AtmNone
		if i == int(d.cursor) {
			focs = style.AtmFocus
		}

		padd := text.EmptyFragment().
			AddSpec(style.SpecPaddingLeft(2))
		indx := d.makeIndex(i, winsize.Cols(digits))
		spac := text.NewFragment(marker.DefaultPaddingText)
		mark := text.NewFragment(o.Text).
			AddAtom(focs)

		lines = append(lines,
			*text.LineFromFragments(*padd, *indx, *spac, *mark),
		)
	}

	unit := drain.UnitFromLines(lines...)
	unit.Drawable.Init()

	d.unit = unit
}

func (d *IndexMenuUnit) makeIndex(cursor int, digits winsize.Cols) *text.Fragment {
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

func (d *IndexMenuUnit) wipe() {
	if d.unit.Drawable.Wipe == nil {
		return
	}
	d.unit.Drawable.Wipe()
}

func (d *IndexMenuUnit) draw(size winsize.Winsize) ([]text.Line, bool) {
	assert.True(d.loaded, drawable.MessageInitialized)

	return d.unit.Drawable.Draw(size)
}
