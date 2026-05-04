package pipeline

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
)

type Transformer func(viewmodel.ViewModel) viewmodel.ViewModel

type Pipeline struct {
	node  screen.Node
	steps []Transformer
}

func New(node screen.Node, steps ...Transformer) *Pipeline {
	return &Pipeline{
		node:  node,
		steps: steps,
	}
}

func (c *Pipeline) PushSteps(steps ...Transformer) *Pipeline {
	c.steps = append(c.steps, steps...)
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

func (c *Pipeline) update(state *state.UIState, event screen.ScreenEvent) screen.Result {
	result := c.node.Screen.Update(state, event)
	if result.Node != nil {
		newNode := New(*result.Node).
			PushSteps(c.steps...).
			ToNode()
		result.Node = &newNode
	}
	return result
}

func (c *Pipeline) view(state state.UIState) viewmodel.ViewModel {
	vm := c.node.Screen.View(state)
	for _, s := range c.steps {
		vm = s(vm)
	}
	return vm
}
