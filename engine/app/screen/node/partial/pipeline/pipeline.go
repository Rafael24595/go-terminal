package pipeline

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
)

type Transformer func(viewmodel.ViewModel) viewmodel.ViewModel

type Pipeline struct {
	node       screen.Node
	steps      []Transformer
	expiration expiration
}

func New(node screen.Node, steps ...Transformer) *Pipeline {
	return &Pipeline{
		node:       node,
		steps:      steps,
		expiration: persistent(),
	}
}

func (c *Pipeline) PushSteps(steps ...Transformer) *Pipeline {
	c.steps = append(c.steps, steps...)
	return c
}

func (c *Pipeline) ExpireOnNode() *Pipeline {
	c.expiration = onNode(&c.node)
	return c
}

func (c *Pipeline) ExpireOnName() *Pipeline {
	c.expiration = onName(c.node.Screen.Name)
	return c
}

func (c *Pipeline) Persistent() *Pipeline {
	c.expiration = persistent()
	return c
}

func (c *Pipeline) ToNode() screen.Node {
	return screen.NewBuilder().
		Name(c.node.Screen.Name).
		AddStack(c.node.Stack).
		Definition(c.node.Screen.Definition).
		Update(c.update).
		View(c.view).
		Children(c.node).
		ToNode()
}

func (c *Pipeline) update(state *state.UIState, event screen.Event) screen.Result {
	result := c.node.Screen.Update(state, event)

	if !c.shouldPropagate(result) {
		return result
	}

	newNode := New(*result.Node).
		PushSteps(c.steps...).
		ToNode()
	result.Node = &newNode

	return result
}

func (c *Pipeline) shouldPropagate(result screen.Result) bool {
	if result.Node == nil {
		return false
	}

	return !c.expiration.on(result.Node)
}

func (c *Pipeline) view(state state.UIState) viewmodel.ViewModel {
	vm := c.node.Screen.View(state)
	for _, s := range c.steps {
		vm = s(vm)
	}
	return vm
}
