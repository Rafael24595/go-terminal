package box

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-terminal/engine/helper/math"
	"github.com/Rafael24595/go-terminal/engine/helper/runes"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable/line"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable/position"
	"github.com/Rafael24595/go-terminal/engine/render/marker"
	"github.com/Rafael24595/go-terminal/engine/render/style"
	"github.com/Rafael24595/go-terminal/engine/render/text"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

const NameBoxDrawable = "BoxDrawable"

const (
	default_padding  = uint(0)
	default_min_size = uint(0)
)

type BoxDrawable struct {
	loaded    bool
	paddingY  uint
	paddingX  uint
	minSize   uint
	textAlign style.HorizontalPosition
	separator marker.BoxSeparatorMeta
	drawable  drawable.Drawable
}

func NewBoxDrawable(drawable drawable.Drawable) *BoxDrawable {
	return &BoxDrawable{
		loaded:    false,
		minSize:   default_min_size,
		paddingY:  default_padding,
		paddingX:  default_padding,
		textAlign: style.Center,
		separator: marker.DefaultBoxSeparator,
		drawable:  drawable,
	}
}

func BoxDrawableFromDrawable(drawable drawable.Drawable) drawable.Drawable {
	return NewBoxDrawable(drawable).ToDrawable()
}

func (d *BoxDrawable) MinSize(size uint) *BoxDrawable {
	d.minSize = size
	return d
}

func (d *BoxDrawable) Separator(separator marker.BoxSeparatorMeta) *BoxDrawable {
	d.separator = separator
	return d
}

func (d *BoxDrawable) PaddingY(padding uint) *BoxDrawable {
	d.paddingY = padding
	return d
}

func (d *BoxDrawable) PaddingX(padding uint) *BoxDrawable {
	d.paddingX = padding
	return d
}

func (d *BoxDrawable) TextAlign(textAlign style.HorizontalPosition) *BoxDrawable {
	d.textAlign = textAlign
	return d
}

func (d *BoxDrawable) ToDrawable() drawable.Drawable {
	return drawable.Drawable{
		Name: NameBoxDrawable,
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

	return position.NewPositionDrawable(d.drawable).
		MarginY(d.paddingY).
		MarginX(d.paddingX).
		PositionY(style.Top).
		PositionX(style.Left).
		ToDrawable()
}

func (d *BoxDrawable) draw(size terminal.Winsize) ([]text.Line, bool) {
	assert.True(d.loaded, "the drawable should be initialized before draw")

	lines, hasNext := d.drawChild(size)

	styled := d.styleLines(size, lines...)

	return styled, hasNext
}

func (d *BoxDrawable) drawChild(size terminal.Winsize) ([]text.Line, bool) {
	lines := make([]text.Line, 0)

	clampSize := d.clampSize(size)

	remaining := int(size.Rows)
	for remaining > 0 {
		line, status := d.drawable.Draw(clampSize)

		if len(line) > 0 {
			lines = append(lines, line...)
		}

		if !status || len(line) == 0 {
			break
		}

		remaining -= len(line)
		if remaining-2 <= 0 {
			assert.Unreachable("box drawables should fit in a single page")
		}
	}

	return lines, remaining <= 0
}

func (d *BoxDrawable) styleLines(size terminal.Winsize, lines ...text.Line) []text.Line {
	vertical := horizontalStaticSize(d.separator)

	minSize := d.minSize + vertical
	maxSize := uint(size.Cols)
	maxLine := drawable.MaxLineSize(lines...)

	padding := math.Clamp(maxLine, minSize, maxSize)

	specCover := style.SpecRepeatLeft(padding + vertical)
	cover := text.LineFromFragments(
		*text.NewFragment(d.separator.Top).AddSpec(specCover),
	)

	result := make([]text.Line, 0)

	result = append(result, *cover)

	available := int(size.Cols) - int(vertical)

	for _, lin := range lines {
		for _, v := range line.WrapLineWords(available, &lin) {
			line := d.styleLine(padding, v)
			result = append(result, line)
		}
	}

	result = append(result, *cover)

	return result
}

func (d *BoxDrawable) styleLine(size uint, line text.Line) text.Line {
	paddingL, paddingR := d.calcPadding(size, line)

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

func (d *BoxDrawable) calcPadding(size uint, line text.Line) (uint, uint) {
	totalWidth := uint(text.LineFragmentsMeasure(&line))

	remaining := size - totalWidth

	switch d.textAlign {
	case style.Left:
		return 0, remaining

	case style.Center:
		paddingL := remaining / 2
		paddingR := remaining - paddingL
		return paddingL, paddingR

	case style.Right:
		return remaining, 0

	}

	assert.Unreachable("undefined justify value: %d", d.textAlign)

	return 0, 0
}

func (d *BoxDrawable) clampSize(size terminal.Winsize) terminal.Winsize {
	vertical := 2
	rows := math.SubClampZero(size.Rows, uint16(vertical))

	horizontal := horizontalStaticSize(d.separator)
	cols := math.SubClampZero(size.Cols, uint16(horizontal))

	return terminal.NewWinsize(rows, cols)
}

func horizontalSeparatorSize(separator marker.BoxSeparatorMeta) (uint, uint) {
	return runes.Measureu(separator.Left), runes.Measureu(separator.Right)
}

func horizontalStaticSize(separator marker.BoxSeparatorMeta) uint {
	left, right := horizontalSeparatorSize(separator)
	return left + right
}
