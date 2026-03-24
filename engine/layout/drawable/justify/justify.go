package justify

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"
	
	"github.com/Rafael24595/go-terminal/engine/layout/drawable"
	"github.com/Rafael24595/go-terminal/engine/render/marker"
	"github.com/Rafael24595/go-terminal/engine/render/style"
	"github.com/Rafael24595/go-terminal/engine/render/text"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

const (
	SlotsBetween = 0
	SlotsEvenly  = 1
	SlotsAround  = 2
)

type JustifyDrawable struct {
	initialized bool
	size        terminal.Winsize
	limit       uint8
	justify     style.Justify
	fragments   []text.Fragment
	cursor      uint16
}

func NewJustifyDrawable(fragments []text.Fragment) *JustifyDrawable {
	return &JustifyDrawable{
		initialized: false,
		size:        terminal.Winsize{},
		limit:       style.DefaultLimit,
		justify:     style.JustifyAround,
		fragments:   fragments,
		cursor:      0,
	}
}

func JustifyDrawableFromFragments(fragments []text.Fragment) drawable.Drawable {
	return NewJustifyDrawable(fragments).ToDrawable()
}

func (d *JustifyDrawable) Limit(limit uint8) *JustifyDrawable {
	d.limit = max(1, limit)
	return d
}

func (d *JustifyDrawable) Justify(justify style.Justify) *JustifyDrawable {
	d.justify = justify
	return d
}

func (d *JustifyDrawable) ToDrawable() drawable.Drawable {
	return drawable.Drawable{
		Init: d.init,
		Draw: d.draw,
	}
}

func (d *JustifyDrawable) init(size terminal.Winsize) {
	d.initialized = true

	d.size = size
}

func (d *JustifyDrawable) draw() ([]text.Line, bool) {
	assert.True(d.initialized, "the drawable should be initialized before draw")

	if d.cursor >= uint16(len(d.fragments)) {
		return make([]text.Line, 0), false
	}

	limit := int(d.limit)
	cols := int(d.size.Cols)

	size := 0
	frags := make([]text.Fragment, 0)

	for i := int(d.cursor); i < len(d.fragments); i++ {
		frag := d.fragments[i]

		fragsLen := len(frags)
		fragSize := text.FragmentMeasure(frag)

		spacing := 0
		if fragsLen > 0 {
			spacing = 1
		}

		newSize := size + spacing + fragSize
		if fragsLen > 0 && fragsLen >= limit || newSize > cols {
			line := justifyLine(cols, frags, size, d.justify)
			return []text.Line{line}, true
		}

		size = newSize
		frags = append(frags, frag)

		d.cursor += 1
	}

	line := justifyLine(cols, frags, size, d.justify)
	return []text.Line{line}, d.cursor < uint16(len(d.fragments))
}

func justifyLine(cols int, frags []text.Fragment, size int, mode style.Justify) text.Line {
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
		next = append(next, space)
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
				text.NewFragment(marker.DefaultPaddingText),
			)
		}
	}
	return spaced
}
