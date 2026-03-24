package inline

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"
	
	"github.com/Rafael24595/go-terminal/engine/layout/drawable"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable/line"
	"github.com/Rafael24595/go-terminal/engine/render/text"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

const NameInlineDrawable = "InlineDrawable"

type InlineDrawable struct {
	initialized bool
	separator   string
	drawables   []drawable.Drawable
	drawable    drawable.Drawable
}

func NewInlineDrawable(drawables ...drawable.Drawable) *InlineDrawable {
	return &InlineDrawable{
		initialized: false,
		separator:   "",
		drawables:   drawables,
		drawable:    drawable.Drawable{},
	}
}

func InlineDrawableFromDrawables(drawables ...drawable.Drawable) drawable.Drawable {
	return NewInlineDrawable(drawables...).ToDrawable()
}

func (d *InlineDrawable) Separator(separator string) *InlineDrawable {
	d.separator = separator
	return d
}

func (d *InlineDrawable) ToDrawable() drawable.Drawable {
	return drawable.Drawable{
		Name: NameInlineDrawable,
		Init: d.init,
		Draw: d.draw,
	}
}

func (d *InlineDrawable) init(size terminal.Winsize) {
	d.initialized = true

	lines := d.drawChildren(size)
	join := d.joinChildren(lines)

	d.drawable = line.EagerDrawableFromLines(join...)

	d.drawable.Init(size)
}

func (d *InlineDrawable) draw() ([]text.Line, bool) {
	assert.True(d.initialized, "the drawable should be initialized before draw")

	return d.drawable.Draw()
}

func (d *InlineDrawable) drawChildren(size terminal.Winsize) []text.Line {
	lines := make([]text.Line, 0)

	if len(d.drawables) == 0 {
		return lines
	}

	index := 0

	focus := d.drawables[index]
	focus.Init(size)

	for {
		result, status := focus.Draw()
		if len(result) > 0 {
			lines = append(lines, result...)
		}

		if status && len(result) == 0 {
			continue
		}

		index += 1
		if index >= len(d.drawables) {
			break
		}

		focus = d.drawables[index]
		focus.Init(size)
	}

	return lines
}

func (d *InlineDrawable) joinChildren(lines []text.Line) []text.Line {
	if len(lines) == 0 {
		return []text.Line{}
	}

	merged := text.EmptyLine()

	var separator *text.Fragment
	if d.separator != "" {
		frag := text.NewFragment(d.separator)
		separator = &frag
	}

	for i, line := range lines {
		frags := line.Text
		if d.separator != "" && i < len(lines)-1 {
			frags = append(frags, *separator)
		}

		merged = merged.SetOrder(line.Order).
			AddFragments(frags...).
			AddSpec(line.Spec)
	}

	return []text.Line{merged}
}
