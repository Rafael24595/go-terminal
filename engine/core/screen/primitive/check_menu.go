package primitive

import (
	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/core/input"
	"github.com/Rafael24595/go-terminal/engine/core/key"
	"github.com/Rafael24595/go-terminal/engine/core/marker"
	"github.com/Rafael24595/go-terminal/engine/core/screen"
	"github.com/Rafael24595/go-terminal/engine/core/style"
	"github.com/Rafael24595/go-terminal/engine/core/text"
	"github.com/Rafael24595/go-terminal/engine/helper/math"
)

const default_check_menu_name = "CheckMenu"

var check_menu_read_definition = screen.DefinitionFromKeys(
	key.NewKeysCode(key.ActionEnter)...,
)

var check_menu_write_definition = screen.DefinitionFromKeys(
	key.NewKeysCode(
		key.ActionEnter,
		key.ActionArrowLeft,
		key.ActionArrowRight,
		key.ActionArrowUp,
		key.ActionArrowDown,
	)...,
)

type CheckMenu struct {
	reference    string
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
	c.limit = min(uint16(len(c.options)), limit)
	return c
}

func (c *CheckMenu) ToScreen() screen.Screen {
	return screen.Screen{
		Name:       c.name,
		Definition: c.definition,
		Update:     c.update,
		View:       c.view,
	}
}

func (c *CheckMenu) definition() screen.Definition {
	return index_menu_definition
}

func (c *CheckMenu) name() string {
	return c.reference
}

func (c *CheckMenu) update(state *state.UIState, event screen.ScreenEvent) screen.ScreenResult {
	//TODO:
	return screen.EmptyScreenResult()
}

func (c *CheckMenu) view(stt state.UIState) core.ViewModel {
	vm := core.ViewModelFromUIState(stt)
	//TODO:
	return *vm
}
