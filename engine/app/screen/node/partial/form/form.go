package form

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/pager"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/config/entry"
	"github.com/Rafael24595/go-reacterm-core/engine/helper/math"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/decorator/inputline"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/drain"
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
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

type Form struct {
	reference string
	items     []entry.Entry
	cursor    uint16
	fixed     bool
}

func New() *Form {
	return &Form{
		reference: Name,
		items:     make([]entry.Entry, 0),
		cursor:    0,
		fixed:     false,
	}
}

func (n *Form) AddNode(
	node screen.Node,
	opts ...entry.Option,
) *Form {
	n.items = append(n.items,
		entry.New(node, opts...),
	)
	return n
}

func (n *Form) ToNode() screen.Node {
	builder := screen.NewBuilder().
		Name(n.reference).
		Definition(n.definition).
		Update(n.update).
		View(n.view)

	for _, v := range n.items {
		builder.Children(v.Node).
			AddStack(v.Node.Stack)
	}

	return builder.ToNode()
}

func (n *Form) definition() screen.Definition {
	local := sources

	item := n.items[n.cursor]
	if item.Selectable {
		local = local.Merge(
			item.Node.Screen.Definition(),
		)
	}

	return local
}

func (n *Form) update(stt *state.UIState, evt screen.Event) screen.Result {
	focus, ok := n.focusItem()

	definition := focus.Node.Screen.Definition()
	required := ok && definition.IsRequired(evt.Key)

	if required {
		result := n.focusUpdate(stt, evt, focus)
		if evt.Key.Code != key.ActionEsc {
			return result
		}
	}

	return n.localUpdate(stt, evt)
}

func (n *Form) localUpdate(stt *state.UIState, evt screen.Event) screen.Result {
	ky := evt.Key

	switch ky.Code {
	case key.ActionEsc:
		n.fixed = false
	case key.ActionArrowUp:
		n.cursor = 0
	case key.ActionArrowDown:
		n.cursor = math.SubClampZeroAs[int, uint16](len(n.items), 1)
	case key.ActionArrowLeft:
		n.cursor = math.SubClampZero(n.cursor, 1)
	case key.ActionArrowRight:
		last := math.SubClampZeroAs[int, uint16](len(n.items), 1)
		n.cursor = min(last, n.cursor+1)
	case key.ActionEnter:
		n.fixed = true
	}

	return screen.ResultFromUIState(stt)
}

func (n *Form) focusUpdate(stt *state.UIState, evt screen.Event, focus entry.Entry) screen.Result {
	result := focus.Node.Screen.Update(stt, evt)

	if result.Node == nil {
		return result

	}

	newItems := make([]entry.Entry, len(n.items))
	copy(newItems, n.items)

	newWrapper := New()
	newWrapper.reference = n.reference
	newWrapper.items = newItems
	newWrapper.cursor = n.cursor
	newWrapper.fixed = n.fixed

	newNode := newWrapper.ToNode()
	result.Node = &newNode

	return result
}

func (n *Form) view(stt state.UIState) viewmodel.ViewModel {
	vm := viewmodel.NewViewModel()

	//TODO: Compile headers and footers?
	for _, i := range n.items {
		cvm := i.Node.Screen.View(stt)

		vm.Kernel.PushLayer(
			cvm.Kernel.ToUnit(),
			i.Opts...
		)

		if cvm.Behavior.NeedsPulse {
			vm.Behavior.NeedsPulse = true
		}
	}

	if focus, ok := n.focusItem(); ok {
		label := focus.Node.Name
		vm.Footer.Push(
			inputline.Wrap(
				drain.UnitFromString(label),
			),
		)
	}

	return *vm
}

func (n *Form) focusItem() (entry.Entry, bool) {
	if n.cursor >= uint16(len(n.items)) {
		return entry.Entry{}, false
	}

	return n.items[n.cursor], true
}
