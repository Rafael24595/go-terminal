package box

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/helper/runes"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/margin"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/utils/drain"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/utils/padding"
	"github.com/Rafael24595/go-reacterm-core/engine/model/hint"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/marker"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
	"github.com/Rafael24595/go-reacterm-core/engine/render/wrap"
)

const Name = "box_unit"

const (
	default_padding = winsize.Cols(0)
)

type BoxUnit struct {
	loaded    bool
	paddingY  winsize.Rows
	paddingX  winsize.Cols
	separator marker.BoxSeparatorMeta
	unit      drawable.Unit
}

func New(unit drawable.Unit) *BoxUnit {
	return &BoxUnit{
		loaded:    false,
		paddingY:  winsize.Rows(default_padding),
		paddingX:  default_padding,
		separator: marker.DefaultBoxSeparator,
		unit:      unit,
	}
}

func Wrap(unit drawable.Unit) drawable.Unit {
	return New(unit).ToUnit()
}

func (d *BoxUnit) Separator(separator marker.BoxSeparatorMeta) *BoxUnit {
	d.separator = separator
	return d
}

func (d *BoxUnit) PaddingY(padding winsize.Rows) *BoxUnit {
	d.paddingY = padding
	return d
}

func (d *BoxUnit) PaddingX(padding winsize.Cols) *BoxUnit {
	d.paddingX = padding
	return d
}

func (d *BoxUnit) ToUnit() drawable.Unit {
	return drawable.NewBuilder().
		Name(Name).
		MergeTags(d.unit.Tags).
		Init(d.init).
		Wipe(d.unit.Drawable.Wipe).
		Draw(d.draw).
		ToUnit()
}

func (d *BoxUnit) init() {
	d.loaded = true

	d.unit = d.makeUnit()

	d.unit.Drawable.Init()
}

func (d *BoxUnit) makeUnit() drawable.Unit {
	if d.paddingY == 0 && d.paddingX == 0 {
		return d.unit
	}

	return margin.NewBuilder().
		Y(d.paddingY, style.Middle).
		X(d.paddingX, style.Center).
		ToUnit(d.unit)
}

func (d *BoxUnit) draw(size winsize.Winsize) ([]text.Line, bool) {
	assert.True(d.loaded, drawable.MessageInitialized)

	innerSize := d.computeInnerSize(size)
	lines, hasNext := drain.UnitLazy(innerSize, d.unit)

	styled := d.styleLines(size, lines...)

	return styled, hasNext
}

// TODO: investigate spec overflow.
func (d *BoxUnit) styleLines(size winsize.Winsize, lines ...text.Line) []text.Line {
	vertical := horizontalStaticSize(d.separator)

	maxLine := text.MaxLineMeasure(size.Cols, lines...)
	measure := min(maxLine+vertical, size.Cols)

	specCover := style.SpecRepeatLeft(measure)
	cover := text.LineFromFragments(
		*text.NewFragment(d.separator.Top).AddSpec(specCover),
	)

	result := make([]text.Line, 0)

	result = append(result, *cover)

	available := size.Cols.Sub(vertical)

	transformer := padding.Cols(
		hint.Fixed(maxLine), style.Center,
	)

	for _, lin := range transformer(size, lines) {
		for _, v := range wrap.Line(available, &lin) {
			line := d.wrapLine(v)
			result = append(result, line)
		}
	}

	result = append(result, *cover)

	return result
}

func (d *BoxUnit) wrapLine(line text.Line) text.Line {
	frags := make([]text.Fragment, 0)

	frags = append(frags, *text.NewFragment(d.separator.Left))
	frags = append(frags, line.Text...)
	frags = append(frags, *text.NewFragment(d.separator.Right))

	line.Text = frags

	return line
}

func (d *BoxUnit) computeInnerSize(size winsize.Winsize) winsize.Winsize {
	vertical := winsize.Rows(2)
	rows := size.Rows.Sub(vertical)

	horizontal := horizontalStaticSize(d.separator)
	cols := size.Cols.Sub(horizontal)

	return winsize.New(rows, cols)
}

func horizontalSeparatorSize(separator marker.BoxSeparatorMeta) (winsize.Cols, winsize.Cols) {
	return runes.Measure(separator.Left), runes.Measure(separator.Right)
}

func horizontalStaticSize(separator marker.BoxSeparatorMeta) winsize.Cols {
	left, right := horizontalSeparatorSize(separator)
	return left + right
}
