package primitive

import (
	"sort"

	"github.com/Rafael24595/go-terminal/engine/app/screen"
	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/app/viewmodel"
	"github.com/Rafael24595/go-terminal/engine/helper/math"
	checkmenu "github.com/Rafael24595/go-terminal/engine/layout/drawable/check"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable/line"
	"github.com/Rafael24595/go-terminal/engine/model/input"
	"github.com/Rafael24595/go-terminal/engine/model/key"
	"github.com/Rafael24595/go-terminal/engine/platform/clock"
	"github.com/Rafael24595/go-terminal/engine/render/marker"
	"github.com/Rafael24595/go-terminal/engine/render/style"
	"github.com/Rafael24595/go-terminal/engine/render/text"
)

const default_check_menu_name = "CheckMenu"

var check_menu_read_definition = screen.DefinitionFromKeys(
	key.NewKeysCode(key.ActionEnter)...,
)

var check_menu_write_definition = screen.DefinitionFromKeys(
	key.NewKeysCode(
		key.ActionEsc,
		key.ActionEnter,
		key.ActionArrowLeft,
		key.ActionArrowRight,
		key.ActionArrowUp,
		key.ActionArrowDown,
	)...,
)

type checkActionHandler = func()

func defaulCheckHandler() {}

type checkAction struct {
	mode    bool
	handler checkActionHandler
}

func enabledCheckAction(handler checkActionHandler) *checkAction {
	return &checkAction{
		mode:    false,
		handler: handler,
	}
}

func disabledCheckAction() *checkAction {
	return &checkAction{
		mode:    false,
		handler: defaulCheckHandler,
	}
}

type CheckMenu struct {
	reference    string
	clock        clock.Clock
	action       *checkAction
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
		action:    disabledCheckAction(),
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

func (c *CheckMenu) definition() screen.Definition {
	if c.action.mode {
		return check_menu_write_definition
	}
	return check_menu_read_definition
}

func (c *CheckMenu) update(state *state.UIState, evnt screen.ScreenEvent) screen.ScreenResult {
	if !c.action.mode {
		return c.updateRead(state, evnt)
	}

	return c.updateNavigation(state, evnt)
}

func (c *CheckMenu) updateNavigation(state *state.UIState, evnt screen.ScreenEvent) screen.ScreenResult {
	ky := evnt.Key

	optsLen := len(c.options)

	switch ky.Code {
	case key.ActionEsc:
		c.action.mode = false
	case key.ActionEnter:
		c.switchState()
	case key.ActionArrowLeft:
		c.cursor = math.SubClampZero(c.cursor, 1)
	case key.ActionArrowRight:
		c.cursor = min(uint(optsLen-1), c.cursor+1)
	case key.ActionArrowUp:
		c.cursor = 0
	case key.ActionArrowDown:
		c.cursor = max(0, uint(optsLen-1))
	}

	return screen.ScreenResultFromUIState(state)
}

func (c *CheckMenu) updateRead(state *state.UIState, evnt screen.ScreenEvent) screen.ScreenResult {
	ky := evnt.Key

	switch ky.Code {
	case key.ActionEnter:
		c.action.mode = true
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

func (c *CheckMenu) view(stt state.UIState) viewmodel.ViewModel {
	indexmenu := checkmenu.NewCheckMenuDrawable(c.options).
		WriteMode(c.action.mode).
		Meta(c.meta).
		Cursor(c.cursor)

	vm := viewmodel.ViewModelFromUIState(stt)

	vm.Header.Shift(
		line.EagerDrawableFromLines(c.title...),
	)
	vm.Lines.Shift(
		indexmenu.ToDrawable(),
	)

	vm.SetStrategy(
		state.NewFocusPager(),
	)

	option := min(len(c.options)-1, int(c.cursor))
	text := c.options[option].Label.Text

	input := viewmodel.NewInputLine(
		line.EagerDrawableFromString(text),
	)

	vm.SetInput(input)

	return *vm
}
