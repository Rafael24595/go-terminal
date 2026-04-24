package primitive

import (
	"sort"

	"github.com/Rafael24595/go-reacterm-core/engine/app/pager"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/set"
	"github.com/Rafael24595/go-reacterm-core/engine/helper/math"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/block"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/widget/checkmenu"
	"github.com/Rafael24595/go-reacterm-core/engine/model/help"
	"github.com/Rafael24595/go-reacterm-core/engine/model/input"
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
	"github.com/Rafael24595/go-reacterm-core/engine/model/param"
	"github.com/Rafael24595/go-reacterm-core/engine/platform/clock"
	"github.com/Rafael24595/go-reacterm-core/engine/render/marker"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

const default_check_menu_name = "CheckMenu"

const ArgCheckMenuActive param.Typed[set.Set[string]] = "check_menu_active"

var check_menu_read_definition = screen.NewDefinitionSources(
	map[key.KeyAction]help.HelpField{
		key.ActionEnter: {Code: []string{"RET"}, Detail: "Edit mode"},
	},
	[]key.KeyAction{
		key.ActionEnter,
	},
)

var check_menu_write_definition = screen.NewDefinitionSources(
	map[key.KeyAction]help.HelpField{
		key.ActionEsc:       {Code: []string{"ESC"}, Detail: "Write Mode"},
		key.ActionEnter:     {Code: []string{"RET"}, Detail: "Active selected"},
		key.ActionArrowUp:   {Code: []string{"↑"}, Detail: "Move first"},
		key.ActionArrowDown: {Code: []string{"↓"}, Detail: "Move last"},
	},
	[]key.KeyAction{
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
	title        []text.Line
	options      []input.CheckOption
	limit        uint16
	cursor       uint
}

func NewCheckMenu() *CheckMenu {
	return &CheckMenu{
		reference: default_check_menu_name,
		clock:     clock.UnixMilliClock,
		action:    input.NewCheckAction(),
		meta:      marker.BracketsCheck,
		title:     make([]text.Line, 0),
		options:   make([]input.CheckOption, 0),
		limit:     0,
		cursor:    0,
	}
}

func (c *CheckMenu) SetName(name string) *CheckMenu {
	c.reference = name
	return c
}

func (c *CheckMenu) SetMeta(meta marker.CheckMeta) *CheckMenu {
	c.meta = meta
	return c
}

func (c *CheckMenu) SetActionHandler(handler input.CheckActionHandler) *CheckMenu {
	c.action.Handler = handler
	return c
}

func (c *CheckMenu) AddTitle(title ...text.Line) *CheckMenu {
	c.title = append(c.title, title...)
	return c
}

func (c *CheckMenu) AddOptions(options ...input.CheckOption) *CheckMenu {
	c.options = append(c.options, options...)
	return c
}

func (c *CheckMenu) SetCursor(cursor uint) *CheckMenu {
	maxIdx := math.SubClampZero(len(c.options), 1)
	c.cursor = math.Clamp(cursor, uint(0), uint(maxIdx))
	return c
}

func (c *CheckMenu) SetDistribution(distribution style.Distribution) *CheckMenu {
	c.distribution = distribution
	return c
}

func (c *CheckMenu) SetLimit(limit uint16) *CheckMenu {
	c.limit = limit
	return c
}

func (c *CheckMenu) ToScreen() screen.Screen {
	screen := screen.Screen{
		Definition: c.definition,
		Update:     c.update,
		View:       c.view,
	}

	return screen.SetName(c.reference).
		StackFromName()
}

func (c *CheckMenu) definitionSource() screen.DefinitionSources {
	if c.action.ActionMode {
		return check_menu_write_definition
	}
	return check_menu_read_definition
}

func (c *CheckMenu) definition() screen.Definition {
	return c.definitionSource().Definition
}

func (c *CheckMenu) update(stt *state.UIState, evt screen.ScreenEvent) screen.ScreenResult {
	if !c.action.ActionMode {
		return c.updateRead(stt, evt)
	}

	return c.updateNavigation(stt, evt)
}

func (c *CheckMenu) updateNavigation(stt *state.UIState, evt screen.ScreenEvent) screen.ScreenResult {
	ky := evt.Key

	optsLen := len(c.options)

	switch ky.Code {
	case key.ActionEsc:
		c.action.ActionMode = false
	case key.ActionEnter:
		c.switchState()
		state.PushParam(
			stt.Stack,
			c.reference,
			ArgCheckMenuActive,
			c.activeIds(),
		)
	case key.ActionArrowLeft:
		c.cursor = math.SubClampZero(c.cursor, 1)
	case key.ActionArrowRight:
		c.cursor = min(uint(optsLen-1), c.cursor+1)
	case key.ActionArrowUp:
		c.cursor = 0
	case key.ActionArrowDown:
		c.cursor = max(0, uint(optsLen-1))
	}

	return screen.ScreenResultFromUIState(stt)
}

func (c *CheckMenu) updateRead(state *state.UIState, evnt screen.ScreenEvent) screen.ScreenResult {
	ky := evnt.Key

	switch ky.Code {
	case key.ActionEnter:
		c.action.ActionMode = true
	}

	return screen.ScreenResultFromUIState(state)
}

func (c *CheckMenu) switchState() {
	optsLen := len(c.options)

	if c.cursor < uint(optsLen) {
		c.options[c.cursor].Status = !c.options[c.cursor].Status
	}

	if c.options[c.cursor].Status {
		c.options[c.cursor].Timestamp = c.clock()
	}

	if c.limit == 0 {
		return
	}

	active := make([]*input.CheckOption, 0, optsLen)
	for i := range c.options {
		if c.options[i].Status {
			active = append(active, &c.options[i])
		}
	}

	if len(active) <= int(c.limit) {
		return
	}

	sort.Slice(active, func(i, j int) bool {
		return active[i].Timestamp < active[j].Timestamp
	})

	excess := len(active) - int(c.limit)
	for i := range excess {
		active[i].Status = false
	}
}

func (c *CheckMenu) activeIds() set.Set[string] {
	result := set.NewSet[string]()
	for _, v := range c.options {
		if v.Status {
			result.Add(v.Id)
		}
	}
	return result
}

func (c *CheckMenu) view(_ state.UIState) viewmodel.ViewModel {
	source := c.definitionSource()

	indexmenu := checkmenu.NewCheckMenuDrawable(c.options).
		WriteMode(c.action.ActionMode).
		Meta(c.meta).
		Cursor(c.cursor)

	vm := viewmodel.NewViewModel()

	vm.Header.Push(
		block.BlockDrawableFromLines(c.title...),
	)
	vm.Kernel.Push(
		indexmenu.ToDrawable(),
	)

	vm.Pager.SetPredicate(
		pager.PredicateFocus(),
	)

	vm.Helper.Push(
		key.ActionsToHelpWithOverride(
			source.Overrides, source.Actions...,
		)...,
	)

	option := min(len(c.options)-1, int(c.cursor))
	text := c.options[option].Label.Text

	input := viewmodel.NewInputLine(
		block.BlockDrawableFromString(text),
	)

	vm.SetInput(input)

	return *vm
}
