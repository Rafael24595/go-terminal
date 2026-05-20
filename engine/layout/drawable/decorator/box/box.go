package box

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/helper/runes"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/spatial/position"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/utils/drain"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/marker"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
	"github.com/Rafael24595/go-reacterm-core/engine/render/wrap"
)

const Name = "box_drawable"

const (
	default_padding  = winsize.Cols(0)
)

type BoxDrawable struct {
	loaded    bool
	paddingY  winsize.Rows
	paddingX  winsize.Cols
	separator marker.BoxSeparatorMeta
	drawable  drawable.Drawable
}

func New(drawable drawable.Drawable) *BoxDrawable {
	return &BoxDrawable{
		loaded:    false,
		paddingY:  winsize.Rows(default_padding),
		paddingX:  default_padding,
		separator: marker.DefaultBoxSeparator,
		drawable:  drawable,
	}
}

func DrawableFromDrawable(drawable drawable.Drawable) drawable.Drawable {
	return New(drawable).ToDrawable()
}

func (d *BoxDrawable) Separator(separator marker.BoxSeparatorMeta) *BoxDrawable {
	d.separator = separator
	return d
}

func (d *BoxDrawable) PaddingY(padding winsize.Rows) *BoxDrawable {
	d.paddingY = padding
	return d
}

func (d *BoxDrawable) PaddingX(padding winsize.Cols) *BoxDrawable {
	d.paddingX = padding
	return d
}

func (d *BoxDrawable) ToDrawable() drawable.Drawable {
	return drawable.Drawable{
		Name: Name,
		Code: d.drawable.Code,
		Tags: d.drawable.Tags,
		Init: d.init,
		Wipe: d.drawable.Wipe,
		Draw: d.draw,
	}
}

func (d *BoxDrawable) init() {
	d.loaded = true

	d.drawable = d.makeDrawable()

	d.drawable.Init()
}

func (d *BoxDrawable) makeDrawable() drawable.Drawable {
	if d.paddingY == 0 && d.paddingX == 0 {
		return d.drawable
	}

	return position.New(d.drawable).
		MarginY(d.paddingY).
		MarginX(d.paddingX).
		PositionY(style.Top).
		PositionX(style.Left).
		ToDrawable()
}

func (d *BoxDrawable) draw(size winsize.Winsize) ([]text.Line, bool) {
	assert.True(d.loaded, drawable.MessageInitialized)

	innerSize := d.computeInnerSize(size)
	lines, hasNext := drain.DrawableLazy(innerSize, d.drawable)

	styled := d.styleLines(size, lines...)

	return styled, hasNext
}

// TODO: investigate spec overflow.
func (d *BoxDrawable) styleLines(size winsize.Winsize, lines ...text.Line) []text.Line {
	vertical := horizontalStaticSize(d.separator)

	maxLine := text.MaxLineMeasure(size.Cols, lines...)
	padding := min(maxLine+vertical, size.Cols)

	specCover := style.SpecRepeatLeft(padding)
	cover := text.LineFromFragments(
		*text.NewFragment(d.separator.Top).AddSpec(specCover),
	)

	result := make([]text.Line, 0)

	result = append(result, *cover)

	available := size.Cols.Sub(vertical)

	for _, lin := range lines {
		for _, v := range wrap.Line(available, &lin) {
			line := d.wrapLine(v)
			result = append(result, line)
		}
	}

	result = append(result, *cover)

	return result
}

func (d *BoxDrawable) wrapLine(line text.Line) text.Line {
	frags := make([]text.Fragment, 0)

	frags = append(frags, *text.NewFragment(d.separator.Left))
	frags = append(frags, line.Text...)
	frags = append(frags, *text.NewFragment(d.separator.Right))

	line.Text = frags

	return line
}

func (d *BoxDrawable) computeInnerSize(size winsize.Winsize) winsize.Winsize {
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
