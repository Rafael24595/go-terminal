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

func (u *CheckMenuUnit) Meta(meta marker.CheckMeta) *CheckMenuUnit {
	u.meta = meta
	return u
}

func (u *CheckMenuUnit) Distribution(distribution style.Distribution) *CheckMenuUnit {
	u.distribution = distribution
	return u
}

func (u *CheckMenuUnit) WriteMode(writeMode bool) *CheckMenuUnit {
	u.writeMode = writeMode
	return u
}

func (u *CheckMenuUnit) Cursor(cursor uint16) *CheckMenuUnit {
	u.cursor = cursor
	return u
}

func (u *CheckMenuUnit) ToUnit() drawable.Unit {
	return drawable.NewBuilder().
		Name(Name).
		Init(u.init).
		Wipe(u.wipe).
		Draw(u.draw).
		ToUnit()
}

func (u *CheckMenuUnit) init() {
	u.initialized = true

	opts := u.addStyles()

	switch u.distribution.Direction {
	case style.Vertical:
		u.unit = u.makeVertical(opts)
	case style.Horizontal:
		u.unit = u.makeHorizontal(opts)
	default:
		assert.Unreachable("undefined direction %d", u.distribution.Direction)

		u.unit = u.makeVertical(opts)
	}

	u.unit.Drawable.Init()
}

func (u *CheckMenuUnit) wipe() {
	if u.unit.Drawable.Wipe == nil {
		return
	}
	u.unit.Drawable.Wipe()
}

func (u *CheckMenuUnit) makeVertical(opts []text.Fragment) drawable.Unit {
	lines := make([]text.Line, len(opts))
	for i := range opts {
		lines[i] = *text.LineFromFragments(opts[i])
	}
	return line.UnitFromLines(lines...)
}

func (u *CheckMenuUnit) makeHorizontal(opts []text.Fragment) drawable.Unit {
	return justify.New(opts).
		Justify(u.distribution.Justify).
		MaxOpts(u.distribution.Limit).
		ToUnit()
}

func (u *CheckMenuUnit) addStyles() []text.Fragment {
	frags := make([]text.Fragment, len(u.options))

	for i := range frags {
		status := u.meta.Unchecked
		if u.options[i].Status {
			status = u.meta.Checked
		}

		label := u.options[i].Label.Text
		if len(label) > 0 {
			label = marker.DefaultPaddingText + label
		}

		frags[i] = *text.NewFragment(u.meta.Open + status + u.meta.Close + label).
			CopyMeta(&u.options[i].Label)

		if u.writeMode && i == int(u.cursor) {
			frags[i].AddAtom(style.AtmSelect, style.AtmFocus)
		}
	}

	return frags
}

func (u *CheckMenuUnit) draw(size winsize.Winsize) ([]text.Line, bool) {
	assert.True(u.initialized, drawable.MessageInitialized)

	return u.unit.Drawable.Draw(size)
}
