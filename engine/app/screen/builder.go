package screen

import (
	"fmt"

	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/set"
	"github.com/Rafael24595/go-reacterm-core/engine/platform/clock"
)

const (
	ErrorMissingName       = "missing_name"
	ErrorMissingDefinition = "missing_definition"
	ErrorMissingUpdate     = "missing_update"
	ErrorMissingView       = "missing_view"
)

func withoutDefinition() Definition {
	return EmptyDefinition()
}

type Builder struct {
	clock      clock.Clock
	name       string
	stack      set.Set[string]
	children   []Node
	definition func() Definition
	update     func(*state.UIState, Event) Result
	view       func(state.UIState) viewmodel.ViewModel
}

func NewBuilder() *Builder {
	return &Builder{
		clock:      clock.GlobalCounterClock,
		name:       "",
		stack:      set.NewSet[string](),
		children:   make([]Node, 0),
		definition: nil,
		update:     nil,
		view:       nil,
	}
}

func (b *Builder) WithClock(clock clock.Clock) *Builder {
	if clock == nil {
		return b
	}

	b.clock = clock
	return b
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

func (b *Builder) Children(children ...Node) *Builder {
	b.children = append(b.children, children...)
	return b
}

func (s *Builder) Update(update UpdateFunc) *Builder {
	s.update = update
	return s
}

func (s *Builder) View(view ViewFunc) *Builder {
	s.view = view
	return s
}

func (s *Builder) makeMeta() Meta {
	meta := NewMeta()

	if s.name == "" {
		meta.Code.Add(ErrorMissingName)
	}

	if s.definition == nil {
		meta.Code.Add(ErrorMissingDefinition)
	}

	if s.update == nil {
		meta.Code.Add(ErrorMissingUpdate)
	}

	if s.view == nil {
		meta.Code.Add(ErrorMissingView)
	}

	return meta
}

func (s *Builder) makeID() string {
	return fmt.Sprintf("%s_%d", s.name, s.clock())
}

func (s *Builder) toScreen() Screen {
	return Screen{
		Name:       s.name,
		Definition: s.definition,
		Update:     s.update,
		View:       s.view,
	}
}

func (s *Builder) ToNode() Node {
	return Node{
		id:       s.makeID(),
		Screen:   s.toScreen(),
		meta:     s.makeMeta(),
		Stack:    s.stack,
		children: s.children,
	}
}
