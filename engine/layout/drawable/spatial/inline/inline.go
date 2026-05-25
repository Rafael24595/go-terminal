package inline

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

const Name = "inline_unit"

type InlineUnit struct {
	loaded    bool
	separator string
	units     []drawable.Unit
}

func New(units ...drawable.Unit) *InlineUnit {
	return &InlineUnit{
		loaded:    false,
		separator: "",
		units:     units,
	}
}

func Wrap(units ...drawable.Unit) drawable.Unit {
	return New(units...).ToUnit()
}

func (u *InlineUnit) Separator(separator string) *InlineUnit {
	u.separator = separator
	return u
}

func (u *InlineUnit) ToUnit() drawable.Unit {
	return drawable.NewBuilder().
		Name(Name).
		Init(u.init).
		Wipe(u.wipe).
		Draw(u.draw).
		ToUnit()
}

func (u *InlineUnit) init() {
	u.loaded = true
}

func (u *InlineUnit) wipe() {}

func (u *InlineUnit) draw(size winsize.Winsize) ([]text.Line, bool) {
	assert.True(u.loaded, drawable.MessageInitialized)

	lines := u.drawChildren(size)

	return u.joinChildren(lines), false
}

func (u *InlineUnit) drawChildren(size winsize.Winsize) []text.Line {
	lines := make([]text.Line, 0)

	if len(u.units) == 0 {
		return lines
	}

	index := 0

	focus := u.units[index]
	focus.Drawable.Init()

	for {
		result, status := focus.Drawable.Draw(size)
		if len(result) > 0 {
			lines = append(lines, result...)
		}

		if status && len(result) == 0 {
			continue
		}

		index += 1
		if index >= len(u.units) {
			break
		}

		focus = u.units[index]
		focus.Drawable.Init()
	}

	return lines
}

func (u *InlineUnit) joinChildren(lines []text.Line) []text.Line {
	if len(lines) == 0 {
		return []text.Line{}
	}

	merged := text.EmptyLine()

	var separator *text.Fragment
	if u.separator != "" {
		frag := text.NewFragment(u.separator)
		separator = frag
	}

	for i, line := range lines {
		frags := line.Text
		if u.separator != "" && i < len(lines)-1 {
			frags = append(frags, *separator)
		}

		merged.PushFragments(frags...).
			CopyMeta(&line)
	}

	return []text.Line{*merged}
}
