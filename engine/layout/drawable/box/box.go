package box

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-terminal/engine/helper/math"
	"github.com/Rafael24595/go-terminal/engine/helper/runes"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable/line"
	"github.com/Rafael24595/go-terminal/engine/render/marker"
	"github.com/Rafael24595/go-terminal/engine/render/style"
	"github.com/Rafael24595/go-terminal/engine/render/text"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

const NameBoxDrawable = "BoxDrawable"

const (
	default_padding  = uint(1)
	default_min_size = uint(0)
)

type BoxDrawable struct {
	initialized bool
	size        terminal.Winsize
	paddingY    uint
	paddingX    uint
	minSize     uint
	positionY   style.VerticalPosition
	positionX   style.HorizontalPosition
	textAlign   style.HorizontalPosition
	spec        style.Spec
	separator   marker.BoxSeparatorMeta
	drawable    drawable.Drawable
}

func NewBoxDrawable(drawable drawable.Drawable) *BoxDrawable {
	return &BoxDrawable{
		initialized: false,
		size:        terminal.Winsize{},
		minSize:     default_min_size,
		paddingY:    default_padding,
		paddingX:    default_padding,
		positionX:   style.Center,
		positionY:   style.Middle,
		textAlign:   style.Center,
		spec:        style.SpecEmpty(),
		separator:   marker.DefaultBoxSeparator,
		drawable:    drawable,
	}
}

func BoxDrawableFromDrawable(drawable drawable.Drawable) drawable.Drawable {
	return NewBoxDrawable(drawable).ToDrawable()
}

func (d *BoxDrawable) MinSize(size uint) *BoxDrawable {
	d.minSize = size
	return d
}

func (d *BoxDrawable) PositionY(vertical style.VerticalPosition) *BoxDrawable {
	d.positionY = vertical
	return d
}

func (d *BoxDrawable) PositionX(horizontal style.HorizontalPosition) *BoxDrawable {
	d.positionX = horizontal
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
		Init: d.init,
		Draw: d.draw,
	}
}

func (d *BoxDrawable) init(size terminal.Winsize) {
	d.initialized = true

	d.spec = makeSpec(d.spec, size, d.positionX)
	d.size = size

	clampSize := d.clampSize(size)
	d.drawable.Init(clampSize)
}

func (d *BoxDrawable) draw() ([]text.Line, bool) {
	assert.True(d.initialized, "the drawable should be initialized before draw")

	lines, hasNext := d.drawChild()

	styled := d.styleLines(lines...)

	base := d.defineBase(styled)
	for _, line := range styled {
		base = append(base, line)
	}

	base = d.fillEmpty(base)
	return base, hasNext
}

func (d *BoxDrawable) defineBase(lines []text.Line) []text.Line {
	size := len(lines)

	if d.positionY == style.Top || size >= int(d.size.Rows) {
		return make([]text.Line, 0)
	}

	start := (d.size.Rows - uint16(size))
	if d.positionY == style.Middle {
		start /= 2
	}

	return make([]text.Line, start)
}

func (d *BoxDrawable) fillEmpty(result []text.Line) []text.Line {
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

func (d *BoxDrawable) drawChild() ([]text.Line, bool) {
	lines := make([]text.Line, 0)

	remaining := int(d.size.Rows)
	for remaining > 0 {
		line, status := d.drawable.Draw()

		if len(line) > 0 {
			lines = append(lines, line...)
		}

		if !status || len(line) == 0 {
			break
		}

		remaining -= len(line)
		if remaining-int(d.paddingY)-2 <= 0 {
			assert.Unreachable("box drawables should fit in a single page")
		}
	}

	return lines, remaining <= 0
}

func (d *BoxDrawable) styleLines(lines ...text.Line) []text.Line {
	vertical := verticalStaticSize(d.separator, d.paddingX)

	minSize := d.minSize + vertical.static
	maxSize := uint(d.size.Cols)
	maxLine := drawable.MaxLineSize(lines...) + vertical.static

	size := math.Clamp(maxLine, minSize, maxSize)

	specCover := style.SpecRepeatLeft(size)
	cover := text.LineFromFragments(
		text.NewFragment(d.separator.Top).AddSpec(specCover),
	).AddSpec(d.spec)

	result := make([]text.Line, 0)

	result = append(result, cover)
	result = d.addPadding(size, vertical.border, result...)

	available := int(d.size.Cols) - int(vertical.static)

	for _, lin := range lines {
		for _, v := range line.WrapLineWords(available, lin) {
			line := d.styleLine(size, vertical, v)
			result = append(result, line)
		}
	}

	result = d.addPadding(size, vertical.border, result...)
	result = append(result, cover)

	return result
}

func (d *BoxDrawable) styleLine(size uint, vertical verticalMeta, line text.Line) text.Line {
	paddingL, paddingR := d.calcPadding(size, vertical, line)

	left := []text.Fragment{
		text.NewFragment(d.separator.Left),
	}

	if paddingL > 0 {
		specLeft := style.SpecRepeatRight(paddingL)
		left = append(left,
			text.NewFragment(d.separator.Space).AddSpec(specLeft),
		)
	}

	right := []text.Fragment{}

	if paddingR > 0 {
		specRight := style.SpecRepeatRight(paddingR)
		right = append(right,
			text.NewFragment(d.separator.Space).AddSpec(specRight),
		)
	}

	right = append(right,
		text.NewFragment(d.separator.Right),
	)

	frags := make([]text.Fragment, 0)

	frags = append(frags, left...)
	frags = append(frags, line.Text...)
	frags = append(frags, right...)

	line.Text = frags
	line.Spec = d.spec

	return line
}

func (d *BoxDrawable) calcPadding(size uint, vertical verticalMeta, line text.Line) (uint, uint) {
	totalWidth := uint(text.LineFragmentsMeasure(line))

	remaining := size - totalWidth - vertical.static
	padding := vertical.padding

	switch d.textAlign {
	case style.Left:
		return padding, remaining + padding

	case style.Center:
		paddingL := remaining / 2
		paddingR := remaining - paddingL
		return paddingL + padding, paddingR + padding

	case style.Right:
		return remaining + padding, padding

	}

	assert.Unreachable("undefined justify value: %d", d.textAlign)

	return 0, 0
}

func (d *BoxDrawable) addPadding(size, borderSize uint, lines ...text.Line) []text.Line {
	available := math.SubClampZero(size, borderSize)

	specSpace := style.SpecRepeatRight(available)
	for range d.paddingY {
		lines = append(lines,
			text.LineFromFragments(
				text.NewFragment(d.separator.Left),
				text.NewFragment(d.separator.Space).AddSpec(specSpace),
				text.NewFragment(d.separator.Right),
			).AddSpec(d.spec),
		)
	}

	return lines
}

func (d *BoxDrawable) clampSize(size terminal.Winsize) terminal.Winsize {
	horizontal := (d.paddingY * 2) + 2
	rows := math.SubClampZero(size.Rows, uint16(horizontal))

	vertical := verticalStaticSize(d.separator, d.paddingX)
	cols := math.SubClampZero(size.Cols, uint16(vertical.static))

	return terminal.NewWinsize(rows, cols)
}

func verticalSeparatorSize(separator marker.BoxSeparatorMeta) (uint, uint) {
	return runes.Measureu(separator.Left), runes.Measureu(separator.Right)
}

func verticalPaddingSize(separator marker.BoxSeparatorMeta, padding uint) uint {
	spaceSize := runes.Measure(separator.Space)
	return padding * uint(spaceSize)
}

type verticalMeta struct {
	static  uint
	border  uint
	padding uint
}

func verticalStaticSize(separator marker.BoxSeparatorMeta, padding uint) verticalMeta {
	left, right := verticalSeparatorSize(separator)
	spaces := verticalPaddingSize(separator, padding)

	boder := left + right

	return verticalMeta{
		static:  boder + (spaces * 2),
		border:  boder,
		padding: spaces,
	}
}

func makeSpec(base style.Spec, size terminal.Winsize, position style.HorizontalPosition) style.Spec {
	cols := uint(size.Cols)

	var spec style.Spec
	switch position {
	case style.Left:
		spec = style.SpecPaddingLeft(cols)
	case style.Center:
		spec = style.SpecPaddingCenter(cols)
	case style.Right:
		spec = style.SpecPaddingRight(cols)
	default:
		return base
	}

	return style.MergeSpec(base, spec)
}
