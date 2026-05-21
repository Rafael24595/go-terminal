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

func (c *Help) ToNode() screen.Node {
	return screen.NewBuilder().
		Name(c.node.Name).
		AddStack(c.node.Stack).
		Definition(c.node.Screen.Definition).
		Update(c.update).
		View(c.view).
		Children(c.node).
		ToNode()
}

func (c *Help) update(state *state.UIState, event screen.Event) screen.Result {
	definition := c.node.Screen.Definition()

	if !definition.IsRequired(event.Key) {
		if event.Key.Code == key.CustomActionHelp {
			c.visible = !c.visible
		}

		state.Helper.ShowHelp = c.visible
		return screen.ResultFromUIState(state)
	}

	c.visible = state.Helper.ShowHelp

	result := c.node.Screen.Update(state, event)
	if result.Node == nil {
		return result
	}

	newWrapper := New(*result.Node)
	newWrapper.visible = c.visible
	newScreen := newWrapper.ToNode()
	result.Node = &newScreen

	return result
}

func (c *Help) view(state state.UIState) viewmodel.ViewModel {
	vm := c.node.Screen.View(state)
	if !c.visible {
		return vm
	}

	definition := c.node.Screen.Definition()

	vm.Footer.Push(
		help.UnitFromFields(definition.Descriptor.ToValuesSlice()),
	)

	return vm
}
