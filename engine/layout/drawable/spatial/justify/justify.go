package justify

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-terminal/engine/commons/structure/set"
	"github.com/Rafael24595/go-terminal/engine/helper/math"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable"
	"github.com/Rafael24595/go-terminal/engine/model/winsize"
	"github.com/Rafael24595/go-terminal/engine/render/marker"
	"github.com/Rafael24595/go-terminal/engine/render/style"
	"github.com/Rafael24595/go-terminal/engine/render/text"
)

const (
	SlotsBetween = 0
	SlotsEvenly  = 1
	SlotsAround  = 2
)

const NameJustifyDrawable = "JustifyDrawable"

type JustifyDrawable struct {
	loaded    bool
	maxOpts   uint16
	maxCols   uint16
	justify   style.Justify
	fragments []text.Fragment
	cursor    uint16
}

func NewJustifyDrawable(fragments []text.Fragment) *JustifyDrawable {
	return &JustifyDrawable{
		loaded:    false,
		maxOpts:   style.DefaultMaxOpts,
		justify:   style.JustifyAround,
		fragments: fragments,
		cursor:    0,
	}
}

func JustifyDrawableFromFragments(fragments []text.Fragment) drawable.Drawable {
	return NewJustifyDrawable(fragments).ToDrawable()
}

func (d *JustifyDrawable) MaxOpts(maxOpts uint16) *JustifyDrawable {
	d.maxOpts = max(1, maxOpts)
	return d
}

func (d *JustifyDrawable) MaxCols(maxCols uint16) *JustifyDrawable {
	d.maxCols = max(1, maxCols)
	return d
}

func (d *JustifyDrawable) Justify(justify style.Justify) *JustifyDrawable {
	d.justify = justify
	return d
}

func (d *JustifyDrawable) AddFragments(fragments []text.Fragment) *JustifyDrawable {
	d.fragments = append(d.fragments, fragments...)
	return d
}

func (d *JustifyDrawable) ToDrawable() drawable.Drawable {
	return drawable.Drawable{
		Name: NameJustifyDrawable,
		Code: "",
		Tags: make(set.Set[string]),
		Init: d.init,
		Draw: d.draw,
		Wipe: d.wipe,
	}
}

func (d *JustifyDrawable) init() {
	d.loaded = true

	d.cursor = 0
}

func (d *JustifyDrawable) wipe() {
	d.cursor = 0
}

func (d *JustifyDrawable) draw(size winsize.Winsize) ([]text.Line, bool) {
	assert.True(d.loaded, drawable.MessageInitialized)

	if d.cursor >= uint16(len(d.fragments)) {
		return make([]text.Line, 0), false
	}

	maxOpts := int(d.maxOpts)
	maxCols := int(math.MinNotZero(size.Cols, d.maxCols))

	width := 0
	frags := make([]text.Fragment, 0)

	for i := int(d.cursor); i < len(d.fragments); i++ {
		frag := d.fragments[i]

		fragsLen := len(frags)
		fragSize := text.FragmentMeasure(int(size.Cols), frag)

		spacing := 0
		if fragsLen > 0 {
			spacing = 1
		}

		newWidth := width + spacing + fragSize
		if fragsLen > 0 && fragsLen >= maxOpts || newWidth > maxCols {
			line := justifyLine(maxCols, frags, width, d.justify)
			return []text.Line{*line}, true
		}

		width = newWidth
		frags = append(frags, frag)

		d.cursor += 1
	}

	line := justifyLine(maxCols, frags, width, d.justify)
	return []text.Line{*line}, d.cursor < uint16(len(d.fragments))
}

func justifyLine(cols int, frags []text.Fragment, size int, mode style.Justify) *text.Line {
	line := text.LineFromFragments(
		addGaps(cols, frags, size, mode)...,
	)

	switch mode {

	case style.JustifyStart:
		return line.AddSpec(style.SpecFromKind(style.SpcKindPaddingRight))

	case style.JustifyEnd:
		return line.AddSpec(style.SpecFromKind(style.SpcKindPaddingLeft))

	case style.JustifyCenter, style.JustifyAround, style.JustifyEvenly:
		return line.AddSpec(style.SpecFromKind(style.SpcKindPaddingCenter))
	}

	return line
}

func addGaps(cols int, frags []text.Fragment, size int, mode style.Justify) []text.Fragment {
	if len(frags) == 0 {
		return frags
	}

	out := make([]text.Fragment, len(frags))
	copy(out, frags)

	free := cols - size
	gaps := len(out) - 1

	if free <= 0 || gaps <= 0 {
		return addSpaceBetween(out)
	}

	switch mode {

	case style.JustifyBetween:
		return distributeSpace(free, out, SlotsBetween)

	case style.JustifyAround:
		return distributeSpace(free, out, SlotsAround)

	case style.JustifyEvenly:
		return distributeSpace(free, out, SlotsEvenly)
	}

	return addSpaceBetween(out)
}

func distributeSpace(free int, frags []text.Fragment, extraSlots int) []text.Fragment {
	gaps := len(frags) - 1

	slots := gaps + extraSlots
	base := free / slots
	remainder := free % slots

	out := make([]text.Fragment, len(frags))
	copy(out, frags)

	fix := 0
	for i := range gaps {
		gap := base
		if remainder > 0 {
			gap++
			remainder--
		}

		if gap <= 0 {
			continue
		}

		space := text.EmptyFragment().AddSpec(
			style.SpecPaddingRight(uint(gap), marker.DefaultPaddingText),
		)

		at := i + fix + 1

		next := make([]text.Fragment, 0, len(out)+1)

		next = append(next, out[:at]...)
		next = append(next, *space)
		next = append(next, out[at:]...)

		out = next

		fix += 1
	}

	return out
}

func addSpaceBetween(frags []text.Fragment) []text.Fragment {
	spaced := make([]text.Fragment, 0, (len(frags)*2)-1)
	for i, f := range frags {
		spaced = append(spaced, f)
		if i < len(frags)-1 {
			spaced = append(spaced,
				*text.NewFragment(marker.DefaultPaddingText),
			)
		}
	}
	return spaced
}
