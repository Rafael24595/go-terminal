package textarea

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/primitive/line"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/widget/textarea/selection"
	"github.com/Rafael24595/go-reacterm-core/engine/model/input"
	"github.com/Rafael24595/go-reacterm-core/engine/model/offset"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/marker"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
	"github.com/Rafael24595/go-reacterm-core/engine/render/wrap"
)

const Name = "text_area_drawable"

type TextAreaDrawable struct {
	loaded     bool
	lazyLoaded bool
	writeMode  bool
	indexMode  bool
	buffer     []rune
	caret      *input.TextCursor
	steps      []Transformer
	drawable   drawable.Drawable
}

func New(buffer []rune, caret *input.TextCursor) *TextAreaDrawable {
	clone := make([]rune, len(buffer))
	copy(clone, buffer)

	return &TextAreaDrawable{
		loaded:     false,
		lazyLoaded: false,
		writeMode:  false,
		indexMode:  false,
		buffer:     clone,
		caret:      caret,
		steps:      make([]Transformer, 0),
		drawable:   drawable.Drawable{},
	}
}

func (d *TextAreaDrawable) WriteMode(writeMode bool) *TextAreaDrawable {
	d.writeMode = writeMode
	return d
}

func (d *TextAreaDrawable) IndexMode(indexMode bool) *TextAreaDrawable {
	d.indexMode = indexMode
	return d
}

func (d *TextAreaDrawable) PushStep(step Transformer) *TextAreaDrawable {
	d.steps = append(d.steps, step)
	return d
}

func (d *TextAreaDrawable) ToDrawable() drawable.Drawable {
	return drawable.Drawable{
		Name: Name,
		Code: d.drawable.Code,
		Tags: d.drawable.Tags,
		Init: d.init,
		Wipe: d.wipe,
		Draw: d.draw,
	}
}

func (d *TextAreaDrawable) init() {
	d.loaded = true
	d.lazyLoaded = false
}

func (d *TextAreaDrawable) lazyInit(size winsize.Winsize) {
	if d.lazyLoaded {
		return
	}

	d.lazyLoaded = true

	start := d.caret.SelectStart().Sub(1)
	end := d.caret.SelectEnd()

	if len(d.buffer) == 0 {
		d.buffer = append(d.buffer, marker.PrintableCaretRunes...)
		start = 0
		end = 1
	}

	frags := d.resolveFragments(d.buffer, start, end)
	for _, step := range d.steps {
		frags = step(frags)
	}

	base := text.LineFromFragments(frags...)

	lines := d.makeLines(base)
	lines = wrap.MaterializeEmpty(size, marker.DefaultPaddingText, lines...)

	drawable := line.DrawableFromLayout(lines...)
	drawable.Init()

	d.drawable = drawable
}

func (d *TextAreaDrawable) makeLines(base *text.Line) []wrap.LayoutLine {
	if d.indexMode {
		return wrap.NormalizeLinesWithOrder(*base)
	}
	return wrap.NormalizeLines(*base)

}

func (d *TextAreaDrawable) wipe() {
	d.lazyLoaded = false

	if d.drawable.Wipe == nil {
		return
	}

	d.drawable.Wipe()
}

func (d *TextAreaDrawable) resolveFragments(
	renderBuffer []rune,
	start, end offset.Offset,
) []text.Fragment {
	frags := make([]text.Fragment, 0, 6)

	if start > 0 {
		frags = append(frags,
			*text.NewFragment(string(renderBuffer[:start])),
		)
	}

	renderer := selection.NewRenderer(
		renderBuffer, start, end, d.blinkStyle(),
	)

	result := renderer.Resolve(d.caret)

	end = result.End
	frags = append(frags, result.Frags...)

	if int(end) < len(renderBuffer) {
		frags = append(frags,
			*text.NewFragment(string(renderBuffer[end:])),
		)
	}

	return frags
}

func (c *TextAreaDrawable) blinkStyle() style.Atom {
	if !c.writeMode {
		return style.AtmNone
	}

	return c.caret.BlinkStyle()
}

func (d *TextAreaDrawable) draw(size winsize.Winsize) ([]text.Line, bool) {
	assert.True(d.loaded, drawable.MessageInitialized)

	d.lazyInit(size)

	return d.drawable.Draw(size)
}
