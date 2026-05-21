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

const Name = "text_area_unit"

type TextAreaUnit struct {
	loaded     bool
	lazyLoaded bool
	writeMode  bool
	indexMode  bool
	buffer     []rune
	caret      *input.TextCursor
	steps      []Transformer
	unit       drawable.Unit
}

func New(buffer []rune, caret *input.TextCursor) *TextAreaUnit {
	clone := make([]rune, len(buffer))
	copy(clone, buffer)

	return &TextAreaUnit{
		loaded:     false,
		lazyLoaded: false,
		writeMode:  false,
		indexMode:  false,
		buffer:     clone,
		caret:      caret,
		steps:      make([]Transformer, 0),
		unit:       drawable.Unit{},
	}
}

func (d *TextAreaUnit) WriteMode(writeMode bool) *TextAreaUnit {
	d.writeMode = writeMode
	return d
}

func (d *TextAreaUnit) IndexMode(indexMode bool) *TextAreaUnit {
	d.indexMode = indexMode
	return d
}

func (d *TextAreaUnit) PushStep(step Transformer) *TextAreaUnit {
	d.steps = append(d.steps, step)
	return d
}

func (d *TextAreaUnit) ToUnit() drawable.Unit {
	return drawable.NewBuilder().
		Name(Name).
		Init(d.init).
		Wipe(d.wipe).
		Draw(d.draw).
		ToUnit()
}

func (d *TextAreaUnit) init() {
	d.loaded = true
	d.lazyLoaded = false
}

func (d *TextAreaUnit) lazyInit(size winsize.Winsize) {
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

	unit := line.UnitFromLayout(lines...)
	unit.Drawable.Init()

	d.unit = unit
}

func (d *TextAreaUnit) makeLines(base *text.Line) []wrap.LayoutLine {
	if d.indexMode {
		return wrap.NormalizeLinesWithOrder(*base)
	}
	return wrap.NormalizeLines(*base)

}

func (d *TextAreaUnit) wipe() {
	d.lazyLoaded = false

	if d.unit.Drawable.Wipe == nil {
		return
	}

	d.unit.Drawable.Wipe()
}

func (d *TextAreaUnit) resolveFragments(
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

func (c *TextAreaUnit) blinkStyle() style.Atom {
	if !c.writeMode {
		return style.AtmNone
	}

	return c.caret.BlinkStyle()
}

func (d *TextAreaUnit) draw(size winsize.Winsize) ([]text.Line, bool) {
	assert.True(d.loaded, drawable.MessageInitialized)

	d.lazyInit(size)

	return d.unit.Drawable.Draw(size)
}
