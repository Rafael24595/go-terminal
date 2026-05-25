package help

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/widget/help"
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
)

type Help struct {
	visible bool
	node    screen.Node
}

func New(node screen.Node) *Help {
	return &Help{
		visible: false,
		node:    node,
	}
}

func (n *Help) ToNode() screen.Node {
	return screen.NewBuilder().
		Name(n.node.Name).
		AddStack(n.node.Stack).
		Definition(n.node.Screen.Definition).
		Update(n.update).
		View(n.view).
		Children(n.node).
		ToNode()
}

func (n *Help) update(state *state.UIState, event screen.Event) screen.Result {
	definition := n.node.Screen.Definition()

	if !definition.IsRequired(event.Key) {
		if event.Key.Code == key.CustomActionHelp {
			n.visible = !n.visible
		}

		state.Helper.ShowHelp = n.visible
		return screen.ResultFromUIState(state)
	}

	n.visible = state.Helper.ShowHelp

	result := n.node.Screen.Update(state, event)
	if result.Node == nil {
		return result
	}

	newWrapper := New(*result.Node)
	newWrapper.visible = n.visible
	newScreen := newWrapper.ToNode()
	result.Node = &newScreen

	return result
}

func (n *Help) view(state state.UIState) viewmodel.ViewModel {
	vm := n.node.Screen.View(state)
	if !n.visible {
		return vm
	}

	definition := n.node.Screen.Definition()

	vm.Footer.Push(
		help.UnitFromFields(definition.Descriptor.ToValuesSlice()),
	)

	return vm
}
