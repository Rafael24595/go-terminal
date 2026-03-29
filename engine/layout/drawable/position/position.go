package position

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-terminal/engine/helper/math"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable"
	"github.com/Rafael24595/go-terminal/engine/render/style"
	"github.com/Rafael24595/go-terminal/engine/render/text"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

const NamePositionDrawable = "PositionDrawable"

const (
	default_margin = uint(0)
)

type PositionDrawable struct {
	initialized bool
	size        terminal.Winsize
	marginY     uint
	marginX     uint
	positionY   style.VerticalPosition
	positionX   style.HorizontalPosition
	spec        style.Spec
	frag        text.Fragment
	drawable    drawable.Drawable
}

func NewPositionDrawable(drawable drawable.Drawable) *PositionDrawable {
	return &PositionDrawable{
		initialized: false,
		size:        terminal.Winsize{},
		positionX:   style.Center,
		positionY:   style.Middle,
		marginY:     default_margin,
		marginX:     default_margin,
		spec:        style.SpecEmpty(),
		frag:        text.EmptyFragment(),
		drawable:    drawable,
	}
}

func PositionDrawableFromDrawable(drawable drawable.Drawable) drawable.Drawable {
	return NewPositionDrawable(drawable).ToDrawable()
}

func (d *PositionDrawable) MarginY(margin uint) *PositionDrawable {
	d.marginY = margin
	return d
}

func (d *PositionDrawable) MarginX(margin uint) *PositionDrawable {
	d.marginX = margin
	return d
}

func (d *PositionDrawable) PositionY(vertical style.VerticalPosition) *PositionDrawable {
	d.positionY = vertical
	return d
}

func (d *PositionDrawable) PositionX(horizontal style.HorizontalPosition) *PositionDrawable {
	d.positionX = horizontal
	return d
}

func (d *PositionDrawable) ToDrawable() drawable.Drawable {
	return drawable.Drawable{
		Name: NamePositionDrawable,
		Init: d.init,
		Draw: d.draw,
	}
}

func (d *PositionDrawable) init(size terminal.Winsize) {
	d.initialized = true

	d.size = size
	d.spec = makeSpec(d.spec, size, d.positionX)
	d.frag = makeFrag(d.frag, d.marginX)

	fixedSize := terminal.Winsize{
		Rows: math.SubClampZero(size.Rows, uint16(d.marginY)*2),
		Cols: math.SubClampZero(size.Cols, uint16(d.marginX)*2),
	}

	d.drawable.Init(fixedSize)
}

func (d *PositionDrawable) draw() ([]text.Line, bool) {
	assert.True(d.initialized, "the drawable should be initialized before draw")

	lines, hasNext := d.drawable.Draw()

	styled := d.styleLines(lines...)

	base := d.makeTopMargin(styled)
	for _, line := range styled {
		line = line.UnshiftFragments(d.frag).
			PushFragments(d.frag)

		base = append(base, line)
	}
	base = d.addBottomMargin(base)

	base = d.fillEmpty(base)
	return base, hasNext
}

func (d *PositionDrawable) makeTopMargin(lines []text.Line) []text.Line {
	size := len(lines)

	if d.positionY == style.Top || size >= int(d.size.Rows) {
		return make([]text.Line, d.marginY)
	}

	start := (d.size.Rows - uint16(size))
	if d.positionY == style.Middle {
		start /= 2
	}

	return make([]text.Line, start + uint16(d.marginY))
}

func (d *PositionDrawable) addBottomMargin(lines []text.Line) []text.Line {
	return append(lines, make([]text.Line, d.marginY)...)
}

func (d *PositionDrawable) fillEmpty(result []text.Line) []text.Line {
	for i := range result {
		if len(result[i].Text) > 0 {
			continue
		}

		result[i].Text = append(
			result[i].Text,
			text.EmptyFragment(),
		)
	}
	return result
}

func (d *PositionDrawable) styleLines(lines ...text.Line) []text.Line {
	for i, v := range lines {
		lines[i] = v.AddSpec(d.spec)
	}
	return lines
}

func makeSpec(base style.Spec, size terminal.Winsize, position style.HorizontalPosition) style.Spec {
	cols := uint(size.Cols)

	var spec style.Spec
	switch position {
	case style.Left:
		spec = style.SpecPaddingRight(cols)
	case style.Center:
		spec = style.SpecPaddingCenter(cols)
	case style.Right:
		spec = style.SpecPaddingLeft(cols)
	default:
		return base
	}

	return style.MergeSpec(base, spec)
}

func makeFrag(frag text.Fragment, margin uint) text.Fragment {
	if margin == 0 {
		return frag
	}

	return text.FragmentFrom(" ", frag).
		AddSpec(style.SpecRepeatRight(margin))
}
