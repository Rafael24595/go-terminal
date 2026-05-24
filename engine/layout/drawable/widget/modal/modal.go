package modal

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/decorator/box"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/spatial/justify"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/spatial/stack"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/drain"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/padding"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/utils/padding/options"
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

func (d *ModalUnit) AddText(text ...text.Line) *ModalUnit {
	d.text = append(d.text, text...)
	return d
}

func (d *ModalUnit) AddOptions(options ...text.Fragment) *ModalUnit {
	d.options = append(d.options, options...)
	return d
}

func (d *ModalUnit) DefineLimit(limit uint) *ModalUnit {
	d.limit = limit
	return d
}

func (d *ModalUnit) DefineCursor(cursor uint16) *ModalUnit {
	d.cursor = cursor
	return d
}

func (d *ModalUnit) ToUnit() drawable.Unit {
	return drawable.NewBuilder().
		Name(Name).
		Init(d.init).
		Wipe(d.wipe).
		Draw(d.draw).
		ToUnit()
}

func (d *ModalUnit) init() {
	d.loaded = true
	d.lazyLoaded = false
}

func (d *ModalUnit) lazyInit(size winsize.Winsize) {
	if d.lazyLoaded {
		return
	}

	d.lazyLoaded = true

	opts := make([]text.Fragment, len(d.options))
	for i := range d.options {
		old := d.options[i]
		opts[i] = *text.NewFragment(old.Text).
			AddAtom(old.Atom).
			AddSpec(old.Spec)

		if i == int(d.cursor) {
			opts[i].AddAtom(style.AtmSelect)
		}
	}

	cols := text.MaxLineMeasure(size.Cols, d.text...) + 1
	text := formatLines(d.text...)

	title := drain.UnitFromLines(text...)

	optionsBlock := drain.Unit(
		justify.New(opts).
			MaxCols(cols).
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
		Y(hint.Maximize[winsize.Rows](), options.WithPosition(style.Middle)).
		X(hint.Maximize[winsize.Cols](), style.Center).
		ToUnit(box)

	position.Drawable.Init()

	d.unit = position
}

func (d *ModalUnit) wipe() {
	d.lazyLoaded = false

	if d.unit.Drawable.Wipe == nil {
		return
	}

	d.unit.Drawable.Wipe()
}

func (d *ModalUnit) draw(size winsize.Winsize) ([]text.Line, bool) {
	assert.True(d.loaded, drawable.MessageInitialized)

	d.lazyInit(size)

	return d.unit.Drawable.Draw(size)
}

func formatLines(lines ...text.Line) []text.Line {
	out := make([]text.Line, len(lines))
	copy(out, lines)

	out = append(out, *text.EmptyLine())

	return out
}
