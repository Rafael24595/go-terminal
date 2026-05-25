package modal

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/config/padding/cols"
	"github.com/Rafael24595/go-reacterm-core/engine/config/padding/rows"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/decorator/box"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/spatial/justify"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/spatial/stack"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/drain"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/padding"
	"github.com/Rafael24595/go-reacterm-core/engine/model/hint"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

const Name = "modal_unit"

type ModalUnit struct {
	loaded     bool
	lazyLoaded bool
	text       []text.Line
	options    []text.Fragment
	limit      uint
	cursor     uint16
	unit       drawable.Unit
}

func New() *ModalUnit {
	return &ModalUnit{
		loaded:     false,
		lazyLoaded: false,
		text:       make([]text.Line, 0),
		options:    make([]text.Fragment, 0),
		limit:      style.DefaultMaxOpts,
		cursor:     0,
		unit:       drawable.Unit{},
	}
}

func (u *ModalUnit) AddText(text ...text.Line) *ModalUnit {
	u.text = append(u.text, text...)
	return u
}

func (u *ModalUnit) AddOptions(options ...text.Fragment) *ModalUnit {
	u.options = append(u.options, options...)
	return u
}

func (u *ModalUnit) DefineLimit(limit uint) *ModalUnit {
	u.limit = limit
	return u
}

func (u *ModalUnit) DefineCursor(cursor uint16) *ModalUnit {
	u.cursor = cursor
	return u
}

func (u *ModalUnit) ToUnit() drawable.Unit {
	return drawable.NewBuilder().
		Name(Name).
		Init(u.init).
		Wipe(u.wipe).
		Draw(u.draw).
		ToUnit()
}

func (u *ModalUnit) init() {
	u.loaded = true
	u.lazyLoaded = false
}

func (u *ModalUnit) lazyInit(size winsize.Winsize) {
	if u.lazyLoaded {
		return
	}

	u.lazyLoaded = true

	opts := make([]text.Fragment, len(u.options))
	for i := range u.options {
		old := u.options[i]
		opts[i] = *text.NewFragment(old.Text).
			AddAtom(old.Atom).
			AddSpec(old.Spec)

		if i == int(u.cursor) {
			opts[i].AddAtom(style.AtmSelect)
		}
	}

	measure := text.MaxLineMeasure(size.Cols, u.text...) + 1
	text := formatLines(u.text...)

	title := drain.UnitFromLines(text...)

	optionsBlock := drain.Unit(
		justify.New(opts).
			MaxCols(measure).
			ToUnit(),
	)

	title.Drawable.Init()
	optionsBlock.Drawable.Init()

	stack := stack.VStackFromUnits(
		title,
		optionsBlock,
	)

	box := box.New(stack).
		PaddingX(1).
		PaddingY(1).
		ToUnit()

	position := padding.NewBuilder().
		Y(hint.Maximize[winsize.Rows](), rows.WithPosition(style.Middle)).
		X(hint.Maximize[winsize.Cols](), cols.WithPosition(style.Center)).
		ToUnit(box)

	position.Drawable.Init()

	u.unit = position
}

func (u *ModalUnit) wipe() {
	u.lazyLoaded = false

	if u.unit.Drawable.Wipe == nil {
		return
	}

	u.unit.Drawable.Wipe()
}

func (u *ModalUnit) draw(size winsize.Winsize) ([]text.Line, bool) {
	assert.True(u.loaded, drawable.MessageInitialized)

	u.lazyInit(size)

	return u.unit.Drawable.Draw(size)
}

func formatLines(lines ...text.Line) []text.Line {
	out := make([]text.Line, len(lines))
	copy(out, lines)

	out = append(out, *text.EmptyLine())

	return out
}
