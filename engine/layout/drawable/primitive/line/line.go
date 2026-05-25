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

func (u *LineUnit) ToUnit() drawable.Unit {
	return drawable.NewBuilder().
		Name(Name).
		Init(u.init).
		Wipe(u.wipe).
		Draw(u.draw).
		ToUnit()
}

func (u *LineUnit) init() {
	u.loaded = true

	u.lines = u.normalizer()
	u.source = wrap.CloneLayoutLines(u.lines...)

	u.indexMeta = computeIndexMeta(u.lines)
}

func (u *LineUnit) wipe() {
	u.source = u.lines
}

func (u *LineUnit) draw(size winsize.Winsize) ([]text.Line, bool) {
	assert.True(u.loaded, drawable.MessageInitialized)

	if len(u.source) == 0 {
		return make([]text.Line, 0), false
	}

	cursor, remain := u.nextIndexedWrappedLine(size)
	u.source = remain

	result := make([]text.Line, 0)
	if cursor != nil {
		result = append(result, *cursor)
	}

	return result, len(u.source) > 0
}

func (u *LineUnit) nextIndexedWrappedLine(size winsize.Winsize) (*text.Line, []wrap.LayoutLine) {
	if u.indexMeta == nil {
		return wrap.NextLine(size.Cols, u.source)
	}
	return NextIndexedLine(size.Cols, u.source, *u.indexMeta)
}
