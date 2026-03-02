package justify

import (
	"github.com/Rafael24595/go-terminal/engine/core/assert"
	"github.com/Rafael24595/go-terminal/engine/core/drawable"
	"github.com/Rafael24595/go-terminal/engine/core/style"
	"github.com/Rafael24595/go-terminal/engine/core/text"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

type JustifyMode uint8

const (
	JustifyStart JustifyMode = iota
	JustifyEnd
	JustifyCenter
	JustifyBetween
	JustifyAround
	JustifyEvenly
)

const (
	SlotsBetween = 0
	SlotsEvenly  = 1
	SlotsAround  = 2
)

const default_limit = 5

type JustifyDrawable struct {
	initialized bool
	size        terminal.Winsize
	limit       uint8
	justify     JustifyMode
	fragments   []text.Fragment
	cursor      uint16
}

func NewJustifyDrawable(fragments []text.Fragment) *JustifyDrawable {
	return &JustifyDrawable{
		initialized: false,
		size:        terminal.Winsize{},
		limit:       default_limit,
		justify:     JustifyAround,
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

func (d *JustifyDrawable) Justify(justify JustifyMode) *JustifyDrawable {
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
		fragsLen := len(frags)

		d.cursor += 1

		if fragsLen > 0 && fragsLen >= limit || size >= cols {
			line := justifyLine(cols, frags, size, d.justify)
			return []text.Line{line}, d.cursor < uint16(len(d.fragments))
		}

		frag := d.fragments[i]

		size += text.FragmentMeasure(frag)
		frags = append(frags, frag)
	}

	line := justifyLine(cols, frags, size, d.justify)
	return []text.Line{line}, d.cursor < uint16(len(d.fragments))
}

func justifyLine(cols int, frags []text.Fragment, size int, mode JustifyMode) text.Line {
	line := text.LineFromFragments(
		addGaps(cols, frags, size, mode)...,
	)

	switch mode {

	case JustifyStart:
		return line.AddSpec(style.SpecFromKind(style.SpcKindPaddingRight))

	case JustifyEnd:
		return line.AddSpec(style.SpecFromKind(style.SpcKindPaddingLeft))

	case JustifyCenter, JustifyAround, JustifyEvenly:
		return line.AddSpec(style.SpecFromKind(style.SpcKindPaddingCenter))
	}

	return line
}

func addGaps(cols int, frags []text.Fragment, size int, mode JustifyMode) []text.Fragment {
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

	case JustifyBetween:
		return distributeSpace(free, out, SlotsBetween)

	case JustifyAround:
		return distributeSpace(free, out, SlotsAround)

	case JustifyEvenly:
		return distributeSpace(free, out, SlotsEvenly)
	}

	return addSpaceBetween(out)
}

func distributeSpace(free int, out []text.Fragment, extraSlots int) []text.Fragment {
	gaps := len(out) - 1

	slots := gaps + extraSlots
	base := free / slots
	remainder := free % slots

	for i := range gaps {
		gap := base
		if remainder > 0 {
			gap++
			remainder--
		}

		if gap <= 0 {
			continue
		}

		out[i].Spec = style.MergeSpec(
			out[i].Spec,
			style.SpecRepeatRight(uint(gap), " "),
		)
	}

	return out
}

func addSpaceBetween(frags []text.Fragment) []text.Fragment {
	spaced := make([]text.Fragment, 0, (len(frags)*2)-1)
	for i, f := range frags {
		spaced = append(spaced, f)
		if i < len(frags)-1 {
			spaced = append(spaced, text.NewFragment(" "))
		}
	}
	return spaced
}
