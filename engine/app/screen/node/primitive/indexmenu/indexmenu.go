package indexmenu

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/pager"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/helper/math"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/decorator/inputline"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/drain"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/widget/indexmenu"
	"github.com/Rafael24595/go-reacterm-core/engine/model/input"
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
	"github.com/Rafael24595/go-reacterm-core/engine/model/param"
	"github.com/Rafael24595/go-reacterm-core/engine/render/marker"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

const Name = "index_menu"

const ArgActiveIndex param.Typed[string] = "id_index_menu"

var index_menu_definition = screen.DefinitionFromActions(
	[]key.Action{
		key.ActionEnter,
		key.ActionArrowLeft,
		key.ActionArrowRight,
		key.ActionArrowUp,
		key.ActionArrowDown,
	}...,
)

type IndexMenu struct {
	reference string
	meta      marker.IndexMeta
	title     []text.Line
	options   []input.MenuOption
	cursor    uint16
}

func New() *IndexMenu {
	return &IndexMenu{
		reference: Name,
		meta:      marker.HyphenIndex,
		title:     make([]text.Line, 0),
		options:   make([]input.MenuOption, 0),
		cursor:    0,
	}
}

func (c *IndexMenu) SetName(name string) *IndexMenu {
	c.reference = name
	return c
}

func (c *IndexMenu) SetMeta(meta marker.IndexMeta) *IndexMenu {
	c.meta = meta
	return c
}

func (c *IndexMenu) AddTitle(title ...text.Line) *IndexMenu {
	c.title = append(c.title, title...)
	return c
}

func (c *IndexMenu) AddOptions(options ...input.MenuOption) *IndexMenu {
	c.options = append(c.options, options...)
	return c
}

func (c *IndexMenu) SetCursor(cursor uint16) *IndexMenu {
	maxIdx := math.SubClampZeroAs[int, uint16](len(c.options), 1)
	c.cursor = math.Clamp(cursor, 0, maxIdx)
	return c
}

func (c *IndexMenu) ToNode() screen.Node {
	return screen.NewBuilder().
		Name(c.reference).
		NameToStack().
		Definition(c.definition).
		Update(c.update).
		View(c.view).
		ToNode()
}

func (c *IndexMenu) definition() screen.Definition {
	return index_menu_definition
}

func (c *IndexMenu) update(stt *state.UIState, evt screen.Event) screen.Result {
	size := uint16(len(c.options))
	if size == 0 {
		return screen.EmptyResult()
	}

	switch evt.Key.Code {
	case key.ActionArrowUp:
		c.cursor = (c.cursor + size - 1) % size
	case key.ActionTab, key.ActionArrowDown:
		c.cursor = (c.cursor + 1) % size
	case key.ActionEnter:
		return c.actionEnter(stt)
	}

	return screen.EmptyResult()
}

func (c *IndexMenu) actionEnter(stt *state.UIState) screen.Result {
	option := c.options[c.cursor]

	state.PushParam(
		stt.Stack,
		c.reference,
		ArgActiveIndex,
		option.Id,
	)

	node := c.options[c.cursor].Action()
	return screen.ResultFromNode(&node)
}

func (c *IndexMenu) view(_ state.UIState) viewmodel.ViewModel {
	frags := input.FragmentFromMenuOption(c.options...)

	indexmenu := indexmenu.New(frags).
		Meta(c.meta).
		Cursor(c.cursor)

	vm := viewmodel.NewViewModel()

	vm.Header.Push(
		drain.DrawableFromLines(c.title...),
	)
	vm.Kernel.Push(
		indexmenu.ToDrawable(),
	)

	index := math.SubClampZeroAs[int, uint16](len(c.options), 1)
	option := min(index, c.cursor)
	text := c.options[option].Label.Text

	vm.Footer.Push(
		inputline.DrawableFromDrawable(
			drain.DrawableFromString(text),
		),
	)

	vm.Pager.SetPredicate(
		pager.PredicateFocus(),
	)

	return *vm
}
