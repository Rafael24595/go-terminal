package line

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
	"github.com/Rafael24595/go-reacterm-core/engine/render/wrap"
)

const Name = "line_unit"

type LineUnit struct {
	loaded     bool
	indexMeta  *indexMeta
	normalizer linesNormalizer
	lines      []wrap.LayoutLine
	source     []wrap.LayoutLine
}

func New(lines ...wrap.LayoutLine) *LineUnit {
	return new(eagerNormalizer(lines...))
}

func FromLines(lines ...text.Line) *LineUnit {
	return new(lazyNormalizer(lines...))
}

func new(normalizer linesNormalizer) *LineUnit {
	return &LineUnit{
		loaded:     false,
		indexMeta:  nil,
		normalizer: normalizer,
	}
}

func UnitFromLayout(lines ...wrap.LayoutLine) drawable.Unit {
	return New(lines...).ToUnit()
}

func UnitFromLines(lines ...text.Line) drawable.Unit {
	return FromLines(lines...).ToUnit()
}

func (d *LineUnit) ToUnit() drawable.Unit {
	return drawable.NewBuilder().
		Name(Name).
		Init(d.init).
		Wipe(d.wipe).
		Draw(d.draw).
		ToUnit()
}

func (d *LineUnit) init() {
	d.loaded = true

	d.lines = d.normalizer()
	d.source = wrap.CloneLayoutLines(d.lines...)

	d.indexMeta = computeIndexMeta(d.lines)
}

func (d *LineUnit) wipe() {
	d.source = d.lines
}

func (d *LineUnit) draw(size winsize.Winsize) ([]text.Line, bool) {
	assert.True(d.loaded, drawable.MessageInitialized)

	if len(d.source) == 0 {
		return make([]text.Line, 0), false
	}

	cursor, remain := d.nextIndexedWrappedLine(size)
	d.source = remain

	result := make([]text.Line, 0)
	if cursor != nil {
		result = append(result, *cursor)
	}

	return result, len(d.source) > 0
}

func (d *LineUnit) nextIndexedWrappedLine(size winsize.Winsize) (*text.Line, []wrap.LayoutLine) {
	if d.indexMeta == nil {
		return wrap.NextLine(size.Cols, d.source)
	}
	return NextIndexedLine(size.Cols, d.source, *d.indexMeta)
}
