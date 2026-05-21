package position

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/marker"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

const Name = "position_unit"

const (
	default_margin = winsize.Cols(0)
)

// TODO: Remove and use padding resources?
type PositionUnit struct {
	loaded    bool
	marginY   winsize.Rows
	marginX   winsize.Cols
	absolute  bool
	positionY style.VerticalPosition
	positionX style.HorizontalPosition
	unit      drawable.Unit
}

func New(unit drawable.Unit) *PositionUnit {
	return &PositionUnit{
		loaded:    false,
		marginY:   winsize.Rows(default_margin),
		marginX:   default_margin,
		absolute:  true,
		positionX: style.Center,
		positionY: style.Middle,
		unit:      unit,
	}
}

func UnitFromUnit(unit drawable.Unit) drawable.Unit {
	return New(unit).ToUnit()
}

func (d *PositionUnit) MarginY(margin winsize.Rows) *PositionUnit {
	d.marginY = margin
	return d
}

func (d *PositionUnit) MarginX(margin winsize.Cols) *PositionUnit {
	d.marginX = margin
	return d
}

func (d *PositionUnit) Absolute(absolute bool) *PositionUnit {
	d.absolute = absolute
	return d
}

func (d *PositionUnit) PositionY(y style.VerticalPosition) *PositionUnit {
	d.positionY = y
	return d
}

func (d *PositionUnit) PositionX(x style.HorizontalPosition) *PositionUnit {
	d.positionX = x
	return d
}

func (d *PositionUnit) ToUnit() drawable.Unit {
	return drawable.NewBuilder().
		Name(Name).
		MergeTags(d.unit.Tags).
		Init(d.init).
		Wipe(d.unit.Drawable.Wipe).
		Draw(d.draw).
		ToUnit()
}

func (d *PositionUnit) init() {
	d.loaded = true

	d.unit.Drawable.Init()
}

func (d *PositionUnit) draw(size winsize.Winsize) ([]text.Line, bool) {
	assert.True(d.loaded, drawable.MessageInitialized)

	fixedSize := winsize.New(
		size.Rows.Sub(d.marginY*2),
		size.Cols.Sub(d.marginX*2),
	)

	spec := makeSpec(fixedSize, d.positionX)
	frag := makeFrag(d.marginX)

	lines, hasNext := d.unit.Drawable.Draw(fixedSize)

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

func (d *PositionUnit) makeTopMargin(size winsize.Winsize, lines []text.Line) []text.Line {
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

func (d *PositionUnit) addBottomMargin(lines []text.Line) []text.Line {
	return append(lines, make([]text.Line, d.marginY)...)
}

func (d *PositionUnit) fillEmpty(result []text.Line) []text.Line {
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

func (d *PositionUnit) styleLines(spec style.Spec, lines ...text.Line) []text.Line {
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
