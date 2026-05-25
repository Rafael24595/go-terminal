package history

import (
	"fmt"

	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/drain"
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

var definition = screen.DefinitionFromActions(
	[]key.Action{
		key.CustomActionBack,
	}...,
)

type History struct {
	history *screen.Node
	node    screen.Node
}

func New(node screen.Node) *History {
	return &History{
		node: node,
	}
}

func (n *History) ToNode() screen.Node {
	return screen.NewBuilder().
		Name(n.node.Name).
		AddStack(n.node.Stack).
		Definition(n.definition).
		Update(n.update).
		View(n.view).
		Children(n.node).
		ToNode()
}

func (n *History) definition() screen.Definition {
	base := n.node.Screen.Definition()
	return definition.Merge(base)
}

func (n *History) update(state *state.UIState, event screen.Event) screen.Result {
	definition := n.node.Screen.Definition()

	if !definition.IsRequired(event.Key) {
		result := n.localUpdate(state, event)
		if result != nil {
			return *result
		}
	}

	result := n.node.Screen.Update(state, event)
	if result.Node == nil {
		return result
	}

	newWrapper := New(*result.Node)
	newWrapper.history = &n.node
	newNode := newWrapper.ToNode()
	result.Node = &newNode

	return result
}

func (n *History) localUpdate(_ *state.UIState, event screen.Event) *screen.Result {
	if n.history == nil || event.Key.Code != key.CustomActionBack {
		return nil
	}

	newBack := New(*n.history)
	newNode := newBack.ToNode()
	result := screen.ResultFromNode(&newNode)

	return &result
}

func (n *History) view(state state.UIState) viewmodel.ViewModel {
	vm := n.node.Screen.View(state)

	if n.history == nil {
		return vm
	}

	page := fmt.Sprintf("back: %s", n.history.Name)

	footer := []text.Line{
		*text.NewLine(page,
			style.SpecFromKind(style.SpcKindPaddingRight),
		),
	}

	vm.Footer.Unshift(
		drain.UnitFromLines(footer...).
			AddTag(screen.SystemMetaTag),
	)

	return vm
}
