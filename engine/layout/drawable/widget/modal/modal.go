package modal

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-terminal/engine/layout/drawable"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable/decorator/box"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable/spatial/justify"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable/spatial/position"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable/spatial/stack"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable/stream/block"
	"github.com/Rafael24595/go-terminal/engine/model/winsize"
	"github.com/Rafael24595/go-terminal/engine/render/style"
	"github.com/Rafael24595/go-terminal/engine/render/text"
)

const NameModalDrawable = "ModalDrawable"

type ModalDrawable struct {
	loaded     bool
	lazyLoaded bool
	text       []text.Line
	options    []text.Fragment
	limit      uint
	cursor     uint
	drawable   drawable.Drawable
}

func NewModalDrawable() *ModalDrawable {
	return &ModalDrawable{
		loaded:     false,
		lazyLoaded: false,
		text:       make([]text.Line, 0),
		options:    make([]text.Fragment, 0),
		limit:      style.DefaultMaxOpts,
		cursor:     0,
		drawable:   drawable.Drawable{},
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
		Code: d.drawable.Code,
		Tags: d.drawable.Tags,
		Init: d.init,
		Wipe: d.wipe,
		Draw: d.draw,
	}
}

func (d *ModalDrawable) init() {
	d.loaded = true
	d.lazyLoaded = false
}

func (d *ModalDrawable) lazyInit(size winsize.Winsize) {
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

	cols := drawable.MaxLineSize(int(size.Cols), d.text...) + 1
	text := formatLines(d.text...)

	title := block.BlockDrawableFromLines(text...)

	options := justify.NewJustifyDrawable(opts).
		MaxCols(uint16(cols)).
		ToDrawable()

	optionsBlock := block.BlockDrawableFromDrawable(options)

	title.Init()
	optionsBlock.Init()

	stack := stack.VStackDrawableFromDrawables(
		title,
		optionsBlock,
	)

	box := box.NewBoxDrawable(stack).
		PaddingX(1).
		PaddingY(1).
		ToDrawable()

	position := position.PositionDrawableFromDrawable(box)
	position.Init()

	d.drawable = position
}

func (d *ModalDrawable) wipe() {
	d.lazyLoaded = false

	if d.drawable.Wipe == nil {
		return
	}

	d.drawable.Wipe()
}

func (d *ModalDrawable) draw(size winsize.Winsize) ([]text.Line, bool) {
	assert.True(d.loaded, drawable.MessageInitialized)

	d.lazyInit(size)

	return d.drawable.Draw(size)
}

func formatLines(lines ...text.Line) []text.Line {
	out := make([]text.Line, len(lines))
	copy(out, lines)

	out = append(out, *text.EmptyLine())

	return out
}
