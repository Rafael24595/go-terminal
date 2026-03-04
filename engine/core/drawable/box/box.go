package box

import (
	"unicode/utf8"

	"github.com/Rafael24595/go-terminal/engine/core/assert"
	"github.com/Rafael24595/go-terminal/engine/core/drawable"
	"github.com/Rafael24595/go-terminal/engine/core/drawable/line"
	"github.com/Rafael24595/go-terminal/engine/core/style"
	"github.com/Rafael24595/go-terminal/engine/core/text"
	"github.com/Rafael24595/go-terminal/engine/helper/math"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

type BoxVerticalPadding uint

const (
	Left BoxVerticalPadding = iota
	Center
	Right
)

type BoxHorizontalPadding uint

const (
	Up BoxHorizontalPadding = iota
	Middle
	Dow
)

const default_inner_padding = uint(1)

type SeparatorMeta struct {
	Top    string
	Bottom string
	Left   string
	Right  string
	Space  string
}

var default_separator = SeparatorMeta{
	Top:    "-",
	Bottom: "-",
	Left:   "|",
	Right:  "|",
	Space:  " ",
}

type BoxDrawable struct {
	initialized  bool
	size         terminal.Winsize
	innerPadding uint
	vertical     BoxVerticalPadding
	horizontal   BoxHorizontalPadding
	spec         style.Spec
	separator    SeparatorMeta
	drawable     drawable.Drawable
}

func NewBoxDrawable(drawable drawable.Drawable) *BoxDrawable {
	return &BoxDrawable{
		initialized:  false,
		size:         terminal.Winsize{},
		innerPadding: default_inner_padding,
		horizontal:   Middle,
		vertical:     Center,
		spec:         style.SpecEmpty(),
		separator:    default_separator,
		drawable:     drawable,
	}
}

func BoxDrawableFromDrawable(drawable drawable.Drawable) drawable.Drawable {
	return NewBoxDrawable(drawable).ToDrawable()
}

func (d *BoxDrawable) Vertical(vertical BoxVerticalPadding) *BoxDrawable {
	d.vertical = vertical
	return d
}

func (d *BoxDrawable) Horizontal(horizontal BoxHorizontalPadding) *BoxDrawable {
	d.horizontal = horizontal
	return d
}

func (d *BoxDrawable) Separator(separator SeparatorMeta) *BoxDrawable {
	d.separator = separator
	return d
}

func (d *BoxDrawable) ToDrawable() drawable.Drawable {
	return drawable.Drawable{
		Init: d.init,
		Draw: d.draw,
	}
}

func (d *BoxDrawable) init(size terminal.Winsize) {
	d.initialized = true

	d.spec = makeSpec(d.spec, size, d.vertical)
	d.size = size

	clampSize := d.clampSize(size)
	d.drawable.Init(clampSize)
}

func (d *BoxDrawable) draw() ([]text.Line, bool) {
	assert.True(d.initialized, "the drawable should be initialized before draw")

	lines, hasNext := d.drawChild()

	styled := d.addStyle(lines...)

	base := d.defineBase(styled)
	for _, line := range styled {
		base = append(base, line)
	}

	base = d.fillEmpty(base)
	return base, hasNext
}

func (d *BoxDrawable) defineBase(lines []text.Line) []text.Line {
	size := len(lines)

	if d.horizontal == Up || size >= int(d.size.Rows) {
		return make([]text.Line, 0)
	}

	start := (d.size.Rows - uint16(size))
	if d.horizontal == Middle {
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
		if remaining-int(d.innerPadding)-2 <= 0 {
			assert.Unreachable("box drawables should fit in a single page")
		}
	}

	return lines, remaining <= 0
}

func (d *BoxDrawable) addStyle(lines ...text.Line) []text.Line {
	borderSize := borderSize(d.separator)

	spaceSize := utf8.RuneCountInString(d.separator.Space)
	paddingSize := d.innerPadding * uint(spaceSize)

	styleSize := uint(borderSize) + (paddingSize * 2)

	size := min(uint(d.size.Cols), drawable.MaxLineSize(lines...)+styleSize)

	specCover := style.SpecRepeatLeft(uint(size))
	cover := text.LineFromFragments(
		text.NewFragment(d.separator.Top).AddSpec(specCover),
	).AddSpec(d.spec)

	result := make([]text.Line, 0)

	result = append(result, cover)
	result = d.addPadding(size, uint(borderSize), result...)

	available := int(d.size.Cols) - int(styleSize)

	for _, lin := range lines {
		for _, v := range line.WrapLineWords(available, lin) {
			totalWidth := uint(text.LineFragmentsMeasure(v))

			leftWidth := uint(utf8.RuneCountInString(d.separator.Left))
			rightWidth := uint(utf8.RuneCountInString(d.separator.Right))

			remaining := size - totalWidth - (leftWidth + rightWidth)

			paddingL := remaining / 2
			paddingR := remaining - paddingL

			specLeft := style.SpecRepeatRight(paddingL)
			specRight := style.SpecRepeatRight(paddingR)

			left := []text.Fragment{
				text.NewFragment(d.separator.Left),
				text.NewFragment(d.separator.Space).AddSpec(specLeft),
			}

			right := []text.Fragment{
				text.NewFragment(d.separator.Space).AddSpec(specRight),
				text.NewFragment(d.separator.Right),
			}

			frags := make([]text.Fragment, 0)

			frags = append(frags, left...)
			frags = append(frags, v.Text...)
			frags = append(frags, right...)

			v.Text = frags
			v.Spec = d.spec

			result = append(result, v)
		}
	}

	result = d.addPadding(size, uint(borderSize), result...)
	result = append(result, cover)

	return result
}

func (d *BoxDrawable) addPadding(size, borderSize uint, lines ...text.Line) []text.Line {
	available := math.SubClampZero(size, borderSize)

	specSpace := style.SpecRepeatRight(available)
	for range d.innerPadding {
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
	horizontal := (d.innerPadding * 2) + 2
	rows := math.SubClampZero(size.Rows, uint16(horizontal))

	vertical := (d.innerPadding * 2) + borderSize(d.separator)
	cols := math.SubClampZero(size.Cols, uint16(vertical))

	return terminal.NewWinsize(rows, cols)
}

func borderSize(separator SeparatorMeta) uint {
	return uint(utf8.RuneCountInString(separator.Left) +
		utf8.RuneCountInString(separator.Right))
}

func makeSpec(base style.Spec, size terminal.Winsize, padding BoxVerticalPadding) style.Spec {
	cols := uint(size.Cols)

	var spec style.Spec
	switch padding {
	case Left:
		spec = style.SpecPaddingLeft(cols)
	case Center:
		spec = style.SpecPaddingCenter(cols)
	case Right:
		spec = style.SpecPaddingRight(cols)
	default:
		return base
	}

	return style.MergeSpec(base, spec)
}
