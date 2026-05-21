package form

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/helper/math"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/decorator/inputline"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/drain"
	"github.com/Rafael24595/go-reacterm-core/engine/model/chunk"
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
)

const Name = "form"

var sources = screen.NewDefinition(
	map[key.Action]key.Descriptor{
		key.ActionEsc:       {Code: []string{"ESC"}, Detail: "Navigation Mode"},
		key.ActionEnter:     {Code: []string{"RET"}, Detail: "Active selected"},
		key.ActionArrowUp:   {Code: []string{"↑"}, Detail: "Move first"},
		key.ActionArrowDown: {Code: []string{"↓"}, Detail: "Move last"},
	},
	[]key.Action{
		key.ActionEsc,
		key.ActionEnter,
		key.ActionArrowLeft,
		key.ActionArrowRight,
		key.ActionArrowUp,
		key.ActionArrowDown,
	},
)

type item struct {
	selectable bool
	node       screen.Node
	chunk      chunk.Chunk[winsize.Rows]
}

type Form struct {
	reference string
	items     []item
	cursor    uint16
	fixed     bool
}

func New() *Form {
	return &Form{
		reference: Name,
		items:     make([]item, 0),
		cursor:    0,
		fixed:     false,
	}
}

func (c *Form) AddNode(
	selectable bool,
	node screen.Node,
	chunk chunk.Chunk[winsize.Rows],
) *Form {
	c.items = append(c.items, item{
		node:       node,
		chunk:      chunk,
		selectable: selectable,
	})
	return c
}

func (c *Form) ToNode() screen.Node {
	builder := screen.NewBuilder().
		Name(c.reference).
		Definition(c.definition).
		Update(c.update).
		View(c.view)

	for _, v := range c.items {
		builder.Children(v.node).
			AddStack(v.node.Stack)
	}

	return builder.ToNode()
}

func (c *Form) definition() screen.Definition {
	local := sources

	item := c.items[c.cursor]
	if item.selectable {
		local = local.Merge(
			item.node.Screen.Definition(),
		)
	}

	return local
}

func (c *Form) update(stt *state.UIState, evt screen.Event) screen.Result {
	focus, ok := c.focusItem()

	definition := focus.node.Screen.Definition()
	required := ok && definition.IsRequired(evt.Key)

	if required {
		result := c.focusUpdate(stt, evt, focus)
		if evt.Key.Code != key.ActionEsc {
			return result
		}
	}

	return c.localUpdate(stt, evt)
}

func (c *Form) localUpdate(stt *state.UIState, evt screen.Event) screen.Result {
	ky := evt.Key

	switch ky.Code {
	case key.ActionEsc:
		c.fixed = false
	case key.ActionArrowUp:
		c.cursor = 0
	case key.ActionArrowDown:
		c.cursor = math.SubClampZeroAs[int, uint16](len(c.items), 1)
	case key.ActionArrowLeft:
		c.cursor = math.SubClampZero(c.cursor, 1)
	case key.ActionArrowRight:
		last := math.SubClampZeroAs[int, uint16](len(c.items), 1)
		c.cursor = min(last, c.cursor+1)
	case key.ActionEnter:
		c.fixed = true
	}

	return screen.ResultFromUIState(stt)
}

func (c *Form) focusUpdate(stt *state.UIState, evt screen.Event, focus item) screen.Result {
	result := focus.node.Screen.Update(stt, evt)

	if result.Node == nil {
		return result

	}

	newItems := make([]item, len(c.items))
	copy(newItems, c.items)

	newWrapper := New()
	newWrapper.reference = c.reference
	newWrapper.items = newItems
	newWrapper.cursor = c.cursor
	newWrapper.fixed = c.fixed

	newNode := newWrapper.ToNode()
	result.Node = &newNode

	return result
}

func (c *Form) view(stt state.UIState) viewmodel.ViewModel {
	vm := viewmodel.NewViewModel()

	//TODO: Compile headers and footers?
	//TODO: Manage the master paging screen.
	for _, i := range c.items {
		cvm := i.node.Screen.View(stt)
		unit := cvm.Kernel.ToUnit()
		vm.Kernel.PushChunk(unit, i.chunk)
	}

	if focus, ok := c.focusItem(); ok {
		label := focus.node.Name
		vm.Footer.Push(
			inputline.UnitFromUnit(
				drain.UnitFromString(label),
			),
		)
	}

	return *vm
}

func (c *Form) focusItem() (item, bool) {
	if c.cursor >= uint16(len(c.items)) {
		return item{}, false
	}

	return c.items[c.cursor], true
}
