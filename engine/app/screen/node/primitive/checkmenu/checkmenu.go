package checkmenu

import (
	"sort"

	"github.com/Rafael24595/go-reacterm-core/engine/app/pager"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/set"
	"github.com/Rafael24595/go-reacterm-core/engine/helper/math"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/decorator/inputline"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/drain"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/widget/checkmenu"
	"github.com/Rafael24595/go-reacterm-core/engine/model/input"
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
	"github.com/Rafael24595/go-reacterm-core/engine/model/param"
	"github.com/Rafael24595/go-reacterm-core/engine/platform/clock"
	"github.com/Rafael24595/go-reacterm-core/engine/render/marker"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
)

const Name = "check_menu"

const ArgActiveChecks param.Typed[set.Set[string]] = "check_menu_active"

var read_definition = screen.NewDefinition(
	map[key.Action]key.Descriptor{
		key.ActionEnter: {Code: []string{"RET"}, Detail: "Edit mode"},
	},
	[]key.Action{
		key.ActionEnter,
	},
)

var write_definition = screen.NewDefinition(
	map[key.Action]key.Descriptor{
		key.ActionEsc:       {Code: []string{"ESC"}, Detail: "Write Mode"},
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

type CheckMenu struct {
	reference    string
	clock        clock.Clock
	action       *input.CheckAction
	meta         marker.CheckMeta
	distribution style.Distribution
	options      []input.CheckOption
	limit        uint16
	cursor       uint16
}

func New() *CheckMenu {
	return &CheckMenu{
		reference: Name,
		clock:     clock.UnixMilliClock,
		action:    input.EmptyCheckAction(),
		meta:      marker.BracketsCheck,
		options:   make([]input.CheckOption, 0),
		limit:     0,
		cursor:    0,
	}
}

func (n *CheckMenu) Name(name string) *CheckMenu {
	n.reference = name
	return n
}

func (n *CheckMenu) Meta(meta marker.CheckMeta) *CheckMenu {
	n.meta = meta
	return n
}

func (n *CheckMenu) ActionHandler(handler input.CheckActionHandler) *CheckMenu {
	n.action.Handler = handler
	return n
}

func (n *CheckMenu) AddOptions(options ...input.CheckOption) *CheckMenu {
	n.options = append(n.options, options...)
	return n
}

func (n *CheckMenu) Cursor(cursor uint16) *CheckMenu {
	maxIdx := math.SubClampZeroAs[int, uint16](len(n.options), 1)
	n.cursor = math.Clamp(cursor, 0, maxIdx)
	return n
}

func (n *CheckMenu) Distribution(distribution style.Distribution) *CheckMenu {
	n.distribution = distribution
	return n
}

func (n *CheckMenu) Limit(limit uint16) *CheckMenu {
	n.limit = limit
	return n
}

func (n *CheckMenu) ToNode() screen.Node {
	return screen.NewBuilder().
		Name(n.reference).
		NameToStack().
		Definition(n.definition).
		Update(n.update).
		View(n.view).
		ToNode()
}

func (n *CheckMenu) definition() screen.Definition {
	if n.action.ActionMode {
		return write_definition
	}
	return read_definition
}

func (n *CheckMenu) update(stt *state.UIState, evt screen.Event) screen.Result {
	if !n.action.ActionMode {
		return n.updateRead(stt, evt)
	}

	return n.updateNavigation(stt, evt)
}

func (n *CheckMenu) updateNavigation(stt *state.UIState, evt screen.Event) screen.Result {
	ky := evt.Key

	optsLen := uint16(len(n.options))

	switch ky.Code {
	case key.ActionEsc:
		n.action.ActionMode = false
	case key.ActionEnter:
		n.switchState()
		state.PushParam(
			stt.Stack,
			n.reference,
			ArgActiveChecks,
			n.activeIds(),
		)
	case key.ActionArrowLeft:
		n.cursor = math.SubClampZero(n.cursor, 1)
	case key.ActionArrowRight:
		optsLen = math.SubClampZero(optsLen, 1)
		n.cursor = min(optsLen, n.cursor+1)
	case key.ActionArrowUp:
		n.cursor = 0
	case key.ActionArrowDown:
		optsLen = math.SubClampZero(optsLen, 1)
		n.cursor = max(0, optsLen)
	}

	return screen.ResultFromUIState(stt)
}

func (n *CheckMenu) updateRead(state *state.UIState, evnt screen.Event) screen.Result {
	ky := evnt.Key

	switch ky.Code {
	case key.ActionEnter:
		n.action.ActionMode = true
	}

	return screen.ResultFromUIState(state)
}

func (n *CheckMenu) switchState() {
	optsLen := uint16(len(n.options))

	if n.cursor < optsLen {
		n.options[n.cursor].Status = !n.options[n.cursor].Status
	}

	if n.options[n.cursor].Status {
		n.options[n.cursor].Timestamp = n.clock()
	}

	if n.limit == 0 {
		return
	}

	active := make([]*input.CheckOption, 0, optsLen)
	for i := range n.options {
		if n.options[i].Status {
			active = append(active, &n.options[i])
		}
	}

	if len(active) <= int(n.limit) {
		return
	}

	sort.Slice(active, func(i, j int) bool {
		return active[i].Timestamp < active[j].Timestamp
	})

	excess := len(active) - int(n.limit)
	for i := range excess {
		active[i].Status = false
	}
}

func (n *CheckMenu) activeIds() set.Set[string] {
	result := set.NewSet[string]()
	for _, v := range n.options {
		if v.Status {
			result.Add(v.Id)
		}
	}
	return result
}

func (n *CheckMenu) view(_ state.UIState) viewmodel.ViewModel {
	indexmenu := checkmenu.New(n.options).
		WriteMode(n.action.ActionMode).
		Meta(n.meta).
		Cursor(n.cursor)

	vm := viewmodel.NewViewModel()

	vm.Kernel.Push(
		indexmenu.ToUnit(),
	)

	vm.Pager.SetPredicate(
		pager.PredicateFocus(),
	)

	index := math.SubClampZeroAs[int, uint16](len(n.options), 1)
	option := min(index, n.cursor)
	text := n.options[option].Label.Text

	vm.Footer.Push(
		inputline.Wrap(
			drain.UnitFromString(text),
		),
	)

	return *vm
}
