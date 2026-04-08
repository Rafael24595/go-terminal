package inline

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-terminal/engine/layout/drawable"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable/block"
	"github.com/Rafael24595/go-terminal/engine/render/text"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

const NameInlineDrawable = "InlineDrawable"

type InlineDrawable struct {
	loaded     bool
	lazyLoaded bool
	size       terminal.Winsize
	separator  string
	drawables  []drawable.Drawable
	drawable   drawable.Drawable
}

func NewInlineDrawable(drawables ...drawable.Drawable) *InlineDrawable {
	return &InlineDrawable{
		loaded:    false,
		size:      terminal.Winsize{},
		separator: "",
		drawables: drawables,
		drawable:  drawable.Drawable{},
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
		Code: d.drawable.Code,
		Tags: d.drawable.Tags,
		Init: d.init,
		Draw: d.draw,
		Wipe: d.wipe,
	}
}

func (d *InlineDrawable) init() {
	d.loaded = true
}

func (d *InlineDrawable) wipe() {
	d.lazyLoaded = false
}

func (d *InlineDrawable) lazyInit(size terminal.Winsize) {
	if d.lazyLoaded {
		return
	}

	d.lazyLoaded = true
	d.size = size

	lines := d.drawChildren()
	join := d.joinChildren(lines)

	d.drawable = block.BlockDrawableFromLines(join...)

	d.drawable.Init()
}

func (d *InlineDrawable) draw(size terminal.Winsize) ([]text.Line, bool) {
	assert.True(d.loaded, "the drawable should be initialized before draw")

	d.lazyInit(size)

	return d.drawable.Draw(size)
}

func (d *InlineDrawable) drawChildren() []text.Line {
	lines := make([]text.Line, 0)

	if len(d.drawables) == 0 {
		return lines
	}

	index := 0

	focus := d.drawables[index]
	focus.Init()

	for {
		result, status := focus.Draw(d.size)
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
		focus.Init()
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

		merged.SetOrder(line.Order).
			PushFragments(frags...).
			AddSpec(line.Spec)
	}

	return []text.Line{merged}
}
