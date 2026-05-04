package screen

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/set"
)

func withoutDefinition() Definition {
	return DefinitionFromKeys()
}

type Builder struct {
	name       string
	stack      set.Set[string]
	meta       ScreenMeta
	definition func() Definition
	update     func(*state.UIState, ScreenEvent) Result
	view       func(state.UIState) viewmodel.ViewModel
}

func NewBuilder() *Builder {
	return &Builder{
		name:       "",
		stack:      set.NewSet[string](),
		meta:       newMeta(),
		definition: nil,
		update:     nil,
		view:       nil,
	}
}

func (s *Builder) Name(name string) *Builder {
	s.name = name
	return s
}

func (s *Builder) NameToStack() *Builder {
	return s.AddStack(
		set.SetFrom(s.name),
	)
}

func (s *Builder) AddStack(stack set.Set[string]) *Builder {
	s.stack.Merge(stack)
	return s
}

func (s *Builder) Definition(definition DefinitionFunc) *Builder {
	s.definition = definition
	return s
}

func (s *Builder) WithoutDefinition() *Builder {
	s.definition = withoutDefinition
	return s
}

func (s *Builder) Update(update UpdateFunc) *Builder {
	s.update = update
	return s
}

func (s *Builder) View(view ViewFunc) *Builder {
	s.view = view
	return s
}

func (s *Builder) ToScreen() Screen {
	return Screen{
		Name:       s.name,
		meta:       s.meta,
		Stack:      s.stack,
		Definition: s.definition,
		Update:     s.update,
		View:       s.view,
	}
}
