package box

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/config/padding/cols"
	"github.com/Rafael24595/go-reacterm-core/engine/config/padding/rows"
	"github.com/Rafael24595/go-reacterm-core/engine/helper/runes"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/margin"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/transform/drain"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/transform/padding"
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

func (u *BoxUnit) Separator(separator marker.BoxSeparatorMeta) *BoxUnit {
	u.separator = separator
	return u
}

func (u *BoxUnit) PaddingY(padding winsize.Rows) *BoxUnit {
	u.paddingY = padding
	return u
}

func (u *BoxUnit) PaddingX(padding winsize.Cols) *BoxUnit {
	u.paddingX = padding
	return u
}

func (u *BoxUnit) ToUnit() drawable.Unit {
	return drawable.NewBuilder().
		Name(Name).
		MergeTags(u.unit.Tags).
		Init(u.init).
		Wipe(u.unit.Drawable.Wipe).
		Draw(u.draw).
		ToUnit()
}

func (u *BoxUnit) init() {
	u.loaded = true

	u.unit = u.makeUnit()

	u.unit.Drawable.Init()
}

func (u *BoxUnit) makeUnit() drawable.Unit {
	if u.paddingY == 0 && u.paddingX == 0 {
		return u.unit
	}

	return margin.NewBuilder().
		Y(u.paddingY, rows.WithPosition(style.Middle)).
		X(u.paddingX, cols.WithPosition(style.Center)).
		ToUnit(u.unit)
}

func (u *BoxUnit) draw(size winsize.Winsize) ([]text.Line, bool) {
	assert.True(u.loaded, drawable.MessageInitialized)

	innerSize := u.computeInnerSize(size)
	lines, hasNext := drain.UnitLazy(innerSize, u.unit)

	styled := u.styleLines(size, lines...)

	return styled, hasNext
}

// TODO: investigate spec overflow.
func (u *BoxUnit) styleLines(size winsize.Winsize, lines ...text.Line) []text.Line {
	vertical := horizontalStaticSize(u.separator)

	maxLine := text.MaxLineMeasure(size.Cols, lines...)
	measure := min(maxLine+vertical, size.Cols)

	specCover := style.SpecRepeatLeft(measure)
	cover := text.LineFromFragments(
		*text.NewFragment(u.separator.Top).AddSpec(specCover),
	)

	result := make([]text.Line, 0)

	result = append(result, *cover)

	available := size.Cols.Sub(vertical)

	transformer := padding.Cols(
		hint.Fixed(maxLine),
		cols.WithPosition(style.Center),
	)

	for _, lin := range transformer(size, lines) {
		for _, v := range wrap.Line(available, &lin) {
			line := u.wrapLine(v)
			result = append(result, line)
		}
	}

	result = append(result, *cover)

	return result
}

func (u *BoxUnit) wrapLine(line text.Line) text.Line {
	frags := make([]text.Fragment, 0)

	frags = append(frags, *text.NewFragment(u.separator.Left))
	frags = append(frags, line.Text...)
	frags = append(frags, *text.NewFragment(u.separator.Right))

	line.Text = frags

	return line
}

func (u *BoxUnit) computeInnerSize(size winsize.Winsize) winsize.Winsize {
	vertical := winsize.Rows(2)
	rows := size.Rows.Sub(vertical)

	horizontal := horizontalStaticSize(u.separator)
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
