package box

import (
	"unicode/utf8"

	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/core/assert"
	"github.com/Rafael24595/go-terminal/engine/core/drawable/line"
	"github.com/Rafael24595/go-terminal/engine/core/style"
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
}

type BoxDrawable struct {
	initialized  bool
	innerPadding uint
	vertical     BoxVerticalPadding
	horizontal   BoxHorizontalPadding
	spec         style.Spec
	separator    SeparatorMeta
	drawable     core.Drawable
	size         terminal.Winsize
}

func NewBoxDrawable(drawable core.Drawable) *BoxDrawable {
	return &BoxDrawable{
		initialized:  false,
		innerPadding: default_inner_padding,
		horizontal:   Middle,
		vertical:     Center,
		spec:         style.SpecEmpty(),
		separator:    default_separator,
		drawable:     drawable,
		size:         terminal.Winsize{},
	}
}

func BoxDrawableFromDrawable(drawable core.Drawable) core.Drawable {
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

func (d *BoxDrawable) ToDrawable() core.Drawable {
	return core.Drawable{
		Init: d.init,
		Draw: d.draw,
	}
}

func (d *BoxDrawable) init(size terminal.Winsize) {
	d.initialized = true

	d.spec = makeSpec(d.spec, size, d.vertical)
	d.size = size

	d.drawable.Init(size)
}

func (d *BoxDrawable) draw() ([]core.Line, bool) {
	assert.True(d.initialized, "the drawable should be initialized before draw")

	lines, hasNext := d.drawChild()

	styled := d.addStyle(lines...)

	base := d.defineBase(styled)
	for _, line := range lines {
		base = append(base, line)
	}

	base = d.fillEmpty(base)
	return base, hasNext
}

func (d *BoxDrawable) defineBase(lines []core.Line) []core.Line {
	size := len(lines)

	if d.horizontal == Up || size >= int(d.size.Rows) {
		return make([]core.Line, 0)
	}

	start := (d.size.Rows - uint16(size))
	if d.horizontal == Middle {
		start /= 2
	}

	return make([]core.Line, start)
}

func (d *BoxDrawable) fillEmpty(result []core.Line) []core.Line {
	for i := range result {
		if len(result[i].Text) == 0 {
			result[i].Text = append(
				result[i].Text,
				core.EmptyFragment(),
			)
		}
	}
	return result
}

func (d *BoxDrawable) drawChild() ([]core.Line, bool) {
	lines := make([]core.Line, 0)

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

func (d *BoxDrawable) addStyle(lines ...core.Line) []core.Line {
	borderSize := utf8.RuneCountInString(d.separator.Left) +
		utf8.RuneCountInString(d.separator.Right)

	spaceSize := utf8.RuneCountInString(d.separator.Space)
	paddingSize := d.innerPadding * uint(spaceSize)

	styleSize := uint(borderSize) + (paddingSize * 2)

	size := min(uint(d.size.Cols), maxLineSize(lines...)+styleSize)

	specCover := style.SpecRepeatRight(uint(size))
	cover := core.LineFromFragments(
		core.NewFragment(d.separator.Top).AddSpec(specCover),
	)

	specSpace := style.SpecRepeatRight(paddingSize)

	left := []core.Fragment{
		core.NewFragment(d.separator.Left),
		core.NewFragment(d.separator.Space).AddSpec(specSpace),
	}

	right := []core.Fragment{
		core.NewFragment(d.separator.Space).AddSpec(specSpace),
		core.NewFragment(d.separator.Right),
	}

	result := make([]core.Line, 0)

	result = append(result, cover)
	result = d.addPadding(size, uint(borderSize), result...)

	available := int(d.size.Cols) - int(styleSize)

	for _, lin := range lines {
		for _, v := range line.WrapLineWords(available, lin) {
			frags := make([]core.Fragment, 0)

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

func (d *BoxDrawable) addPadding(size, borderSize uint, lines ...core.Line) []core.Line {
	available := math.SubClampZero(size, borderSize)

	specSpace := style.SpecRepeatRight(available)
	for range d.innerPadding {
		lines = append(lines,
			core.LineFromFragments(
				core.NewFragment(d.separator.Left),
				core.NewFragment(d.separator.Space).AddSpec(specSpace),
				core.NewFragment(d.separator.Right),
			),
		)
	}

	return lines
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

func maxLineSize(lines ...core.Line) uint {
	size := uint(0)
	for _, v := range lines {
		measure := core.LineFragmentsMeasure(v)
		size = max(size, uint(measure))
	}
	return size
}
