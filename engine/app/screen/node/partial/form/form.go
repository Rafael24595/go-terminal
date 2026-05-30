package form

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/node/partial/dummy"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/node/partial/pipeline"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/config/entry"
	"github.com/Rafael24595/go-reacterm-core/engine/config/layer"
	"github.com/Rafael24595/go-reacterm-core/engine/helper/math"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/decorator/inputline"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/gutter"
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
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
		key.CustomActionPointer,
	},
)

type Form struct {
	reference string
	pointer   uint8
	focused   bool
	cursor    uint16
	steps     []pipeline.Transformer
	items     []entry.Entry
}

func New() *Form {
	return &Form{
		reference: Name,
		pointer:   0,
		focused:   false,
		cursor:    0,
		steps:     make([]pipeline.Transformer, 0),
		items:     make([]entry.Entry, 0),
	}
}

func (n *Form) PushSteps(steps ...pipeline.Transformer) *Form {
	n.steps = append(n.steps, steps...)
	return n
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

func (n *Form) AddBreak(rows ...winsize.Rows) *Form {
	fixed := winsize.Rows(1)
	if len(rows) > 0 {
		fixed = rows[0]
	}

	return n.AddNode(
		dummy.ToNode(),
		entry.WithLayout(
			layer.Fixed(fixed),
		),
	)
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
		n.focused = false
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
		n.focused = true
	case key.CustomActionPointer:
		n.pointer = nextPointer(n.pointer)
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

	newSteps := make([]pipeline.Transformer, len(n.steps))
	copy(newSteps, n.steps)

	newWrapper := New()
	newWrapper.reference = n.reference
	newWrapper.pointer = n.pointer
	newWrapper.focused = n.focused
	newWrapper.cursor = n.cursor
	newWrapper.steps = newSteps
	newWrapper.items = newItems

	newNode := newWrapper.ToNode()
	result.Node = &newNode

	return result
}

func (n *Form) view(stt state.UIState) viewmodel.ViewModel {
	vm := viewmodel.New()

	pointer := findPointer(n.pointer)

	//TODO: Compile headers and footers?
	for i, e := range n.items {
		cvm := e.Node.Screen.View(stt)

		opts := make([]gutter.Option, 0, 1)

		if pointer.hasNone(pointerGutter) || n.cursor != uint16(i) {
			opts = append(opts,
				gutter.WithLeftGutter(gutter.DefaultEmpty),
			)
		}

		unit := gutter.Unit(
			cvm.Kernel.ToUnit(),
			opts...,
		)

		vm.Kernel.PushLayer(unit, e.Opts...)

		if cvm.Behavior.NeedsPulse {
			vm.Behavior.NeedsPulse = true
		}
	}

	focus, ok := n.focusItem()
	if ok && pointer.hasAny(pointerPrompt) {
		label := text.NewFragment(focus.Node.Name).
			AddAtom(style.AtmSelect)

		vm.Footer.Push(
			inputline.FromFragment(*label),
		)
	}

	return n.applySteps(*vm)
}

func (n *Form) applySteps(vm viewmodel.ViewModel) viewmodel.ViewModel {
	for _, s := range n.steps {
		vm = s(vm)
	}
	return vm
}

func (n *Form) focusItem() (entry.Entry, bool) {
	if n.cursor >= uint16(len(n.items)) {
		return entry.Entry{}, false
	}

	return n.items[n.cursor], true
}
