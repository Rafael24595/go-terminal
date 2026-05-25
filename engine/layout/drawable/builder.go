package drawable

import (
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/set"
)

const (
	ErrorMissingName = "missing_name"
	ErrorMissingInit = "missing_init"
	ErrorMissingWipe = "missing_wipe"
	ErrorMissingDraw = "missing_draw"
)

type Builder struct {
	name string
	tags set.Set[string]
	init InitFunc
	wipe WipeFunc
	draw DrawFunc
}

func NewBuilder() *Builder {
	return &Builder{
		name: "",
		tags: set.NewSet[string](),
		init: nil,
		wipe: nil,
		draw: nil,
	}
}

func (b *Builder) Name(name string) *Builder {
	b.name = name
	return b
}

func (b *Builder) AddTags(tags ...string) *Builder {
	b.tags.Add(tags...)
	return b
}

func (b *Builder) MergeTags(tags set.Set[string]) *Builder {
	b.tags.Merge(tags)
	return b
}

func (b *Builder) Init(init InitFunc) *Builder {
	b.init = init
	return b
}

func (b *Builder) Wipe(wipe WipeFunc) *Builder {
	b.wipe = wipe
	return b
}

func (b *Builder) Draw(draw DrawFunc) *Builder {
	b.draw = draw
	return b
}

func (b *Builder) makeTags() set.Set[string] {
	tags := set.NewSet[string]()

	if b.name == "" {
		tags.Add(ErrorMissingName)
	}

	if b.init == nil {
		b.tags.Add(ErrorMissingInit)
	}

	if b.wipe == nil {
		b.tags.Add(ErrorMissingWipe)
	}

	if b.draw == nil {
		b.tags.Add(ErrorMissingDraw)
	}

	tags.Merge(b.tags)

	return tags
}

func (b *Builder) toDrawable() Drawable {
	return Drawable{
		Init: b.init,
		Wipe: b.wipe,
		Draw: b.draw,
	}
}

func (b *Builder) ToUnit() Unit {
	return Unit{
		Name:     b.name,
		Tags:     b.makeTags(),
		Drawable: b.toDrawable(),
	}
}
