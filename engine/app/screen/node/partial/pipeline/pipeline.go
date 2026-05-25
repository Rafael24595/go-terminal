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

func (n *Pipeline) PushSteps(steps ...Transformer) *Pipeline {
	n.steps = append(n.steps, steps...)
	return n
}

func (n *Pipeline) ExpireOnNode() *Pipeline {
	n.expiration = onNode(&n.node)
	return n
}

func (n *Pipeline) ExpireOnName() *Pipeline {
	n.expiration = onName(n.node.Name)
	return n
}

func (n *Pipeline) Persistent() *Pipeline {
	n.expiration = persistent()
	return n
}

func (n *Pipeline) ToNode() screen.Node {
	return screen.NewBuilder().
		Name(n.node.Name).
		AddStack(n.node.Stack).
		Definition(n.node.Screen.Definition).
		Update(n.update).
		View(n.view).
		Children(n.node).
		ToNode()
}

func (n *Pipeline) update(state *state.UIState, event screen.Event) screen.Result {
	result := n.node.Screen.Update(state, event)

	if !n.shouldPropagate(result) {
		return result
	}

	newNode := New(*result.Node).
		PushSteps(n.steps...).
		ToNode()
	result.Node = &newNode

	return result
}

func (n *Pipeline) shouldPropagate(result screen.Result) bool {
	if result.Node == nil {
		return false
	}

	return !n.expiration.on(result.Node)
}

func (n *Pipeline) view(state state.UIState) viewmodel.ViewModel {
	vm := n.node.Screen.View(state)
	for _, s := range n.steps {
		vm = s(vm)
	}
	return vm
}
