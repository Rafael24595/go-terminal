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

func (b *Builder) Name(name string) *Builder {
	b.name = name
	return b
}

func (b *Builder) NameToStack() *Builder {
	return b.AddStack(
		set.SetFrom(b.name),
	)
}

func (b *Builder) AddStack(stack set.Set[string]) *Builder {
	b.stack.Merge(stack)
	return b
}

func (b *Builder) Definition(definition DefinitionFunc) *Builder {
	b.definition = definition
	return b
}

func (b *Builder) WithoutDefinition() *Builder {
	b.definition = withoutDefinition
	return b
}

func (b *Builder) Children(children ...Node) *Builder {
	b.children = append(b.children, children...)
	return b
}

func (b *Builder) Update(update UpdateFunc) *Builder {
	b.update = update
	return b
}

func (b *Builder) View(view ViewFunc) *Builder {
	b.view = view
	return b
}

func (b *Builder) makeTags() set.Set[string] {
	tags := set.NewSet[string]()

	if b.name == "" {
		tags.Add(ErrorMissingName)
	}

	if b.definition == nil {
		tags.Add(ErrorMissingDefinition)
	}

	if b.update == nil {
		tags.Add(ErrorMissingUpdate)
	}

	if b.view == nil {
		tags.Add(ErrorMissingView)
	}

	return tags
}

func (b *Builder) makeID() string {
	return fmt.Sprintf("%s_%d", b.name, b.clock())
}

func (b *Builder) toScreen() Screen {
	return Screen{
		Definition: b.definition,
		Update:     b.update,
		View:       b.view,
	}
}

func (b *Builder) ToNode() Node {
	return Node{
		id:       b.makeID(),
		Name:     b.name,
		Tags:     b.makeTags(),
		Screen:   b.toScreen(),
		Stack:    b.stack,
		children: b.children,
	}
}
