package screen

import (
	"fmt"

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
	definition DefinitionFunc
	update     UpdateFunc
	view       ViewFunc
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

func (s *Builder) makeTags() set.Set[string] {
	tags := set.NewSet[string]()

	if s.name == "" {
		tags.Add(ErrorMissingName)
	}

	if s.definition == nil {
		tags.Add(ErrorMissingDefinition)
	}

	if s.update == nil {
		tags.Add(ErrorMissingUpdate)
	}

	if s.view == nil {
		tags.Add(ErrorMissingView)
	}

	return tags
}

func (s *Builder) makeID() string {
	return fmt.Sprintf("%s_%d", s.name, s.clock())
}

func (s *Builder) toScreen() Screen {
	return Screen{
		Definition: s.definition,
		Update:     s.update,
		View:       s.view,
	}
}

func (s *Builder) ToNode() Node {
	return Node{
		id:       s.makeID(),
		Name:     s.name,
		Tags:     s.makeTags(),
		Screen:   s.toScreen(),
		Stack:    s.stack,
		children: s.children,
	}
}
