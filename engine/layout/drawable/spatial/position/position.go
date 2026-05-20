package position

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/marker"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

const Name = "position_drawable"

const (
	default_margin = winsize.Cols(0)
)

//TODO: Remove and use padding resources?
type PositionDrawable struct {
	loaded    bool
	marginY   winsize.Rows
	marginX   winsize.Cols
	absolute  bool
	positionY style.VerticalPosition
	positionX style.HorizontalPosition
	drawable  drawable.Drawable
}

func New(drawable drawable.Drawable) *PositionDrawable {
	return &PositionDrawable{
		loaded:    false,
		marginY:   winsize.Rows(default_margin),
		marginX:   default_margin,
		absolute:  true,
		positionX: style.Center,
		positionY: style.Middle,
		drawable:  drawable,
	}
}

func DrawableFromDrawable(drawable drawable.Drawable) drawable.Drawable {
	return New(drawable).ToDrawable()
}

func (d *PositionDrawable) MarginY(margin winsize.Rows) *PositionDrawable {
	d.marginY = margin
	return d
}

func (d *PositionDrawable) MarginX(margin winsize.Cols) *PositionDrawable {
	d.marginX = margin
	return d
}

func (d *PositionDrawable) Absolute(absolute bool) *PositionDrawable {
	d.absolute = absolute
	return d
}

func (d *PositionDrawable) PositionY(y style.VerticalPosition) *PositionDrawable {
	d.positionY = y
	return d
}

func (d *PositionDrawable) PositionX(x style.HorizontalPosition) *PositionDrawable {
	d.positionX = x
	return d
}

func (d *PositionDrawable) ToDrawable() drawable.Drawable {
	return drawable.Drawable{
		Name: Name,
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

func (d *PositionDrawable) draw(size winsize.Winsize) ([]text.Line, bool) {
	assert.True(d.loaded, drawable.MessageInitialized)

	fixedSize := winsize.New(
		size.Rows.Sub(d.marginY*2),
		size.Cols.Sub(d.marginX*2),
	)

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

func (d *PositionDrawable) makeTopMargin(size winsize.Winsize, lines []text.Line) []text.Line {
	width := winsize.Rows(len(lines))

	if d.positionY == style.Top || width >= size.Rows {
		return make([]text.Line, d.marginY)
	}

	start := size.Rows - width
	if d.positionY == style.Middle {
		start /= 2
	}

	return make([]text.Line, start+d.marginY)
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

func makeSpec(size winsize.Winsize, position style.HorizontalPosition) style.Spec {
	switch position {
	case style.Left:
		return style.SpecPaddingRight(size.Cols)
	case style.Center:
		return style.SpecPaddingCenter(size.Cols)
	case style.Right:
		return style.SpecPaddingLeft(size.Cols)
	}

	return style.SpecEmpty()
}

func makeFrag(margin winsize.Cols) *text.Fragment {
	if margin == 0 {
		return text.EmptyFragment()
	}

	return text.NewFragment(marker.DefaultPaddingText).
		AddSpec(style.SpecRepeatRight(margin))
}
