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

func (c *History) ToNode() screen.Node {
	return screen.NewBuilder().
		Name(c.node.Screen.Name).
		AddStack(c.node.Stack).
		Definition(c.definition).
		Update(c.update).
		View(c.view).
		Children(c.node).
		ToNode()
}

func (c *History) definition() screen.Definition {
	base := c.node.Screen.Definition()
	return definition.Merge(base)
}

func (c *History) update(state *state.UIState, event screen.Event) screen.Result {
	definition := c.node.Screen.Definition()

	if !definition.IsRequired(event.Key) {
		result := c.localUpdate(state, event)
		if result != nil {
			return *result
		}
	}

	result := c.node.Screen.Update(state, event)
	if result.Node == nil {
		return result
	}

	newWrapper := New(*result.Node)
	newWrapper.history = &c.node
	newNode := newWrapper.ToNode()
	result.Node = &newNode

	return result
}

func (c *History) localUpdate(_ *state.UIState, event screen.Event) *screen.Result {
	if c.history == nil || event.Key.Code != key.CustomActionBack {
		return nil
	}

	newBack := New(*c.history)
	newNode := newBack.ToNode()
	result := screen.ResultFromNode(&newNode)

	return &result
}

func (c *History) view(state state.UIState) viewmodel.ViewModel {
	vm := c.node.Screen.View(state)

	if c.history == nil {
		return vm
	}

	page := fmt.Sprintf("back: %s", c.history.Screen.Name)

	footer := []text.Line{
		*text.NewLine(page,
			style.SpecFromKind(style.SpcKindPaddingRight),
		),
	}

	vm.Footer.Unshift(
		drain.DrawableFromLines(footer...).
			AddTag(screen.SystemMetaTag),
	)

	return vm
}
