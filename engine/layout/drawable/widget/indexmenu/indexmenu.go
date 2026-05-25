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

func (u *IndexMenuUnit) Meta(meta marker.IndexMeta) *IndexMenuUnit {
	u.meta = meta
	return u
}

func (u *IndexMenuUnit) Cursor(cursor uint16) *IndexMenuUnit {
	u.cursor = cursor
	return u
}

func (u *IndexMenuUnit) ToUnit() drawable.Unit {
	return drawable.NewBuilder().
		Name(Name).
		Init(u.init).
		Wipe(u.wipe).
		Draw(u.draw).
		ToUnit()
}

func (u *IndexMenuUnit) init() {
	u.loaded = true

	lines := make([]text.Line, 0)

	digits := math.Digits(len(u.options))

	for i, o := range u.options {
		focs := style.AtmNone
		if i == int(u.cursor) {
			focs = style.AtmFocus
		}

		padd := text.EmptyFragment().
			AddSpec(style.SpecPaddingLeft(2))
		indx := u.makeIndex(i, winsize.Cols(digits))
		spac := text.NewFragment(marker.DefaultPaddingText)
		mark := text.NewFragment(o.Text).
			AddAtom(focs)

		lines = append(lines,
			*text.LineFromFragments(*padd, *indx, *spac, *mark),
		)
	}

	unit := drain.UnitFromLines(lines...)
	unit.Drawable.Init()

	u.unit = unit
}

func (u *IndexMenuUnit) makeIndex(cursor int, digits winsize.Cols) *text.Fragment {
	if u.meta.Kind == marker.Numeric {
		txt := helper.Right(strconv.Itoa(cursor+1), digits)
		index := text.NewFragment(txt + ".- ")
		if cursor == int(u.cursor) {
			index.Atom |= style.AtmBold
		}
		return index
	}

	if u.meta.Kind == marker.Alphabetic {
		txt := helper.Right(helper.NumberToAlpha(cursor), digits)
		index := text.NewFragment(txt + ".- ")
		if cursor == int(u.cursor) {
			index.Atom |= style.AtmBold
		}
		return index
	}

	index := u.meta.Index
	if cursor == int(u.cursor) {
		index = u.meta.Cursor
	}

	return text.NewFragment(index)
}

func (u *IndexMenuUnit) wipe() {
	if u.unit.Drawable.Wipe == nil {
		return
	}
	u.unit.Drawable.Wipe()
}

func (u *IndexMenuUnit) draw(size winsize.Winsize) ([]text.Line, bool) {
	assert.True(u.loaded, drawable.MessageInitialized)

	return u.unit.Drawable.Draw(size)
}
