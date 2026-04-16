package position

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-terminal/engine/helper/math"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable"
	"github.com/Rafael24595/go-terminal/engine/render/marker"
	"github.com/Rafael24595/go-terminal/engine/render/style"
	"github.com/Rafael24595/go-terminal/engine/render/text"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

const NamePositionDrawable = "PositionDrawable"

const (
	default_margin = uint(0)
)

type PositionDrawable struct {
	loaded    bool
	marginY   uint
	marginX   uint
	absolute  bool
	positionY style.VerticalPosition
	positionX style.HorizontalPosition
	drawable  drawable.Drawable
}

func NewPositionDrawable(drawable drawable.Drawable) *PositionDrawable {
	return &PositionDrawable{
		loaded:    false,
		marginY:   default_margin,
		marginX:   default_margin,
		absolute:  true,
		positionX: style.Center,
		positionY: style.Middle,
		drawable:  drawable,
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

func (d *PositionDrawable) Absolute(absolute bool) *PositionDrawable {
	d.absolute = absolute
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
		Code: d.drawable.Code,
		Tags: d.drawable.Tags,
		Init: d.init,
		Wipe: d.drawable.Init,
		Draw: d.draw,
	}
}

func (d *PositionDrawable) init() {
	d.loaded = true

	d.drawable.Init()
}

func (d *PositionDrawable) draw(size terminal.Winsize) ([]text.Line, bool) {
	assert.True(d.loaded, drawable.MessageInitialized)

	fixedSize := terminal.Winsize{
		Rows: math.SubClampZero(size.Rows, uint16(d.marginY)*2),
		Cols: math.SubClampZero(size.Cols, uint16(d.marginX)*2),
	}

	spec := makeSpec(fixedSize, d.positionX)
	frag := makeFrag(d.marginX)

	lines, hasNext := d.drawable.Draw(fixedSize)

	styled := d.styleLines(spec, lines...)

	base := d.makeTopMargin(size, styled)
	for _, line := range styled {
		line.UnshiftFragments(*frag).
			PushFragments(*frag)

		base = append(base, line)
	}
	base = d.addBottomMargin(base)

	base = d.fillEmpty(base)

	return base, hasNext
}

func (d *PositionDrawable) makeTopMargin(size terminal.Winsize, lines []text.Line) []text.Line {
	width := len(lines)

	if d.positionY == style.Top || width >= int(size.Rows) {
		return make([]text.Line, d.marginY)
	}

	start := (size.Rows - uint16(width))
	if d.positionY == style.Middle {
		start /= 2
	}

	return make([]text.Line, start+uint16(d.marginY))
}

func (d *PositionDrawable) addBottomMargin(lines []text.Line) []text.Line {
	return append(lines, make([]text.Line, d.marginY)...)
}

func (d *PositionDrawable) fillEmpty(result []text.Line) []text.Line {
	var frag text.Fragment
	if d.absolute {
		frag = *text.EmptyFragment()
	} else {
		frag = *text.NewFragment(marker.DefaultPaddingText)
	}

	for i := range result {
		if len(result[i].Text) > 0 {
			continue
		}

		result[i].Text = append(
			result[i].Text,
			frag,
		)
	}
	return result
}

func (d *PositionDrawable) styleLines(spec style.Spec, lines ...text.Line) []text.Line {
	for i := range lines {
		lines[i].AddSpec(spec)
	}
	return lines
}

func makeSpec(size terminal.Winsize, position style.HorizontalPosition) style.Spec {
	cols := uint(size.Cols)

	switch position {
	case style.Left:
		return style.SpecPaddingRight(cols)
	case style.Center:
		return style.SpecPaddingCenter(cols)
	case style.Right:
		return style.SpecPaddingLeft(cols)
	}

	return style.SpecEmpty()
}

func makeFrag(margin uint) *text.Fragment {
	if margin == 0 {
		return text.EmptyFragment()
	}

	return text.NewFragment(marker.DefaultPaddingText).
		AddSpec(style.SpecRepeatRight(margin))
}
