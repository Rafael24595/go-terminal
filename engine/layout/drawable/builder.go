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

func (s *Builder) Name(name string) *Builder {
	s.name = name
	return s
}

func (s *Builder) AddTags(tags ...string) *Builder {
	s.tags.Add(tags...)
	return s
}

func (s *Builder) MergeTags(tags set.Set[string]) *Builder {
	s.tags.Merge(tags)
	return s
}

func (s *Builder) Init(init InitFunc) *Builder {
	s.init = init
	return s
}

func (s *Builder) Wipe(wipe WipeFunc) *Builder {
	s.wipe = wipe
	return s
}

func (s *Builder) Draw(draw DrawFunc) *Builder {
	s.draw = draw
	return s
}

func (s *Builder) makeTags() set.Set[string] {
	tags := set.NewSet[string]()

	if s.name == "" {
		tags.Add(ErrorMissingName)
	}

	if s.init == nil {
		s.tags.Add(ErrorMissingInit)
	}

	if s.wipe == nil {
		s.tags.Add(ErrorMissingWipe)
	}

	if s.draw == nil {
		s.tags.Add(ErrorMissingDraw)
	}

	tags.Merge(s.tags)

	return tags
}

func (s *Builder) toDrawable() Drawable {
	return Drawable{
		Init: s.init,
		Wipe: s.wipe,
		Draw: s.draw,
	}
}

func (s *Builder) ToUnit() Unit {
	return Unit{
		Name:     s.name,
		Tags:     s.makeTags(),
		Drawable: s.toDrawable(),
	}
}
