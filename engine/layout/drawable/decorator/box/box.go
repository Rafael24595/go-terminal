package box

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/helper/math"
	"github.com/Rafael24595/go-reacterm-core/engine/helper/runes"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/primitive/line"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/spatial/position"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/marker"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

const Name = "box_drawable"

const (
	default_padding  = winsize.Cols(0)
	default_min_size = winsize.Cols(0)
)

type BoxDrawable struct {
	loaded    bool
	paddingY  winsize.Rows
	paddingX  winsize.Cols
	minSize   winsize.Cols
	textAlign style.HorizontalPosition
	separator marker.BoxSeparatorMeta
	drawable  drawable.Drawable
}

func New(drawable drawable.Drawable) *BoxDrawable {
	return &BoxDrawable{
		loaded:    false,
		minSize:   default_min_size,
		paddingY:  winsize.Rows(default_padding),
		paddingX:  default_padding,
		textAlign: style.Center,
		separator: marker.DefaultBoxSeparator,
		drawable:  drawable,
	}
}

func DrawableFromDrawable(drawable drawable.Drawable) drawable.Drawable {
	return New(drawable).ToDrawable()
}

func (d *BoxDrawable) MinSize(size winsize.Cols) *BoxDrawable {
	d.minSize = size
	return d
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

func (d *BoxDrawable) TextAlign(textAlign style.HorizontalPosition) *BoxDrawable {
	d.textAlign = textAlign
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

	lines, hasNext := d.drawChild(size)

	styled := d.styleLines(size, lines...)

	return styled, hasNext
}

func (d *BoxDrawable) drawChild(size winsize.Winsize) ([]text.Line, bool) {
	lines := make([]text.Line, 0)

	clampSize := d.clampSize(size)

	remaining := size.Rows
	for remaining > 0 {
		line, status := d.drawable.Draw(clampSize)

		lineLen := winsize.Rows(len(line))
		if lineLen > 0 {
			lines = append(lines, line...)
		}

		if !status || lineLen == 0 {
			break
		}

		remaining = remaining.Clamp(lineLen)
	}

	return lines, remaining <= 0
}

func (d *BoxDrawable) styleLines(size winsize.Winsize, lines ...text.Line) []text.Line {
	vertical := horizontalStaticSize(d.separator)

	minSize := d.minSize + vertical
	maxSize := size.Cols
	maxLine := drawable.MaxLineSize(size.Cols, lines...)

	padding := math.Clamp(maxLine, minSize, maxSize)

	specCover := style.SpecRepeatLeft(padding + vertical)
	cover := text.LineFromFragments(
		*text.NewFragment(d.separator.Top).AddSpec(specCover),
	)

	result := make([]text.Line, 0)

	result = append(result, *cover)

	available := size.Cols.Clamp(vertical)

	for _, lin := range lines {
		for _, v := range line.WrapLineWords(available, &lin) {
			line := d.styleLine(padding, v)
			result = append(result, line)
		}
	}

	result = append(result, *cover)

	return result
}

func (d *BoxDrawable) styleLine(cols winsize.Cols, line text.Line) text.Line {
	paddingL, paddingR := d.calcPadding(cols, line)

	left := []text.Fragment{
		*text.NewFragment(d.separator.Left),
	}

	if paddingL > 0 {
		specLeft := style.SpecRepeatRight(paddingL)
		left = append(left,
			*text.NewFragment(d.separator.Space).AddSpec(specLeft),
		)
	}

	right := []text.Fragment{}

	if paddingR > 0 {
		specRight := style.SpecRepeatRight(paddingR)
		right = append(right,
			*text.NewFragment(d.separator.Space).AddSpec(specRight),
		)
	}

	right = append(right,
		*text.NewFragment(d.separator.Right),
	)

	frags := make([]text.Fragment, 0)

	frags = append(frags, left...)
	frags = append(frags, line.Text...)
	frags = append(frags, right...)

	line.Text = frags

	return line
}

func (d *BoxDrawable) calcPadding(cols winsize.Cols, line text.Line) (winsize.Cols, winsize.Cols) {
	totalWidth := text.FragmentMeasure(cols, line.Text...)

	remaining := cols.Clamp(totalWidth)

	switch d.textAlign {
	case style.Left:
		return 0, remaining

	case style.Center:
		paddingL := remaining / 2
		paddingR := remaining.Clamp(paddingL)
		return paddingL, paddingR

	case style.Right:
		return remaining, 0

	}

	assert.Unreachable("undefined justify value: %d", d.textAlign)

	return 0, 0
}

func (d *BoxDrawable) clampSize(size winsize.Winsize) winsize.Winsize {
	vertical := winsize.Rows(2)
	rows := size.Rows.Clamp(vertical)

	horizontal := horizontalStaticSize(d.separator)
	cols := size.Cols.Clamp(horizontal)

	return winsize.New(rows, cols)
}

func horizontalSeparatorSize(separator marker.BoxSeparatorMeta) (winsize.Cols, winsize.Cols) {
	return runes.Measure(separator.Left), runes.Measure(separator.Right)
}

func horizontalStaticSize(separator marker.BoxSeparatorMeta) winsize.Cols {
	left, right := horizontalSeparatorSize(separator)
	return left + right
}
