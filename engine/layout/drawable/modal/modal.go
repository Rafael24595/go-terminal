package modal

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"
	
	"github.com/Rafael24595/go-terminal/engine/layout/drawable"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable/box"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable/justify"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable/line"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable/stack"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable/static"
	"github.com/Rafael24595/go-terminal/engine/render/style"
	"github.com/Rafael24595/go-terminal/engine/render/text"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

const NameModalDrawable = "ModalDrawable"

type ModalDrawable struct {
	initialized bool
	text        []text.Line
	options     []text.Fragment
	limit       uint
	cursor      uint
	box         drawable.Drawable
}

func NewModalDrawable() *ModalDrawable {
	return &ModalDrawable{
		initialized: false,
		text:        make([]text.Line, 0),
		options:     make([]text.Fragment, 0),
		limit:       style.DefaultLimit,
		cursor:      0,
		box:         drawable.Drawable{},
	}
}

func ModalDrawableFromData(text []text.Line, options []text.Fragment, cursor uint) drawable.Drawable {
	return NewModalDrawable().ToDrawable()
}

func (d *ModalDrawable) AddText(text ...text.Line) *ModalDrawable {
	d.text = append(d.text, text...)
	return d
}

func (d *ModalDrawable) AddOptions(options ...text.Fragment) *ModalDrawable {
	d.options = append(d.options, options...)
	return d
}

func (d *ModalDrawable) DefineLimit(limit uint) *ModalDrawable {
	d.limit = limit
	return d
}

func (d *ModalDrawable) DefineCursor(cursor uint) *ModalDrawable {
	d.cursor = cursor
	return d
}

func (d *ModalDrawable) ToDrawable() drawable.Drawable {
	return drawable.Drawable{
		Name: NameModalDrawable,
		Init: d.init,
		Draw: d.draw,
	}
}

func (d *ModalDrawable) init(size terminal.Winsize) {
	d.initialized = true

	opts := make([]text.Fragment, len(d.options))
	for i := range d.options {
		old := d.options[i]
		opts[i] = text.NewFragment(old.Text).
			AddAtom(old.Atom).
			AddSpec(old.Spec)

		if i == int(d.cursor) {
			opts[i] = opts[i].AddAtom(style.AtmSelect)
		}
	}

	cols := drawable.MaxLineSize(d.text...) + 1
	text := formatLines(d.text...)

	eager := line.EagerDrawableFromLines(text...)
	justify := justify.JustifyDrawableFromFragments(opts)

	eager.Init(size)
	justify.Init(terminal.Winsize{
		Rows: size.Rows,
		Cols: uint16(cols),
	})

	stack := stack.StackDrawableFromDrawables(
		static.StaticDrawableFromDrawable(eager),
		static.StaticDrawableFromDrawable(justify),
	)

	box := box.BoxDrawableFromDrawable(stack)
	box.Init(size)

	d.box = box
}

func (d *ModalDrawable) draw() ([]text.Line, bool) {
	assert.True(d.initialized, "the drawable should be initialized before draw")

	return d.box.Draw()
}

func formatLines(lines ...text.Line) []text.Line {
	out := make([]text.Line, len(lines))
	copy(out, lines)

	out = append(out, text.EmptyLine())

	return out
}
