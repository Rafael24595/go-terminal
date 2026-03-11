package wrapper

import (
	"fmt"

	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/core/drawable/line"
	"github.com/Rafael24595/go-terminal/engine/core/help"
	"github.com/Rafael24595/go-terminal/engine/core/key"
	"github.com/Rafael24595/go-terminal/engine/core/screen"
	"github.com/Rafael24595/go-terminal/engine/core/style"
	"github.com/Rafael24595/go-terminal/engine/core/text"
	"github.com/Rafael24595/go-terminal/engine/helper/math"
)

var pagination_overrides = map[key.KeyAction]help.HelpField{
	key.ActionArrowLeft:  {Code: []string{"←"}, Detail: "Prev page"},
	key.ActionArrowRight: {Code: []string{"→"}, Detail: "Next page"},
}

var pagination_actions = []key.KeyAction{
	key.ActionArrowLeft,
	key.ActionArrowRight,
}

var pagination_keys = key.NewKeysCode(
	pagination_actions...,
)

type Pagination struct {
	screen screen.Screen
}

func NewPagination(screen screen.Screen) *Pagination {
	return &Pagination{
		screen: screen,
	}
}

func (c *Pagination) ToScreen() screen.Screen {
	return screen.Screen{
		Name:       c.screen.Name,
		Definition: c.definition,
		Update:     c.update,
		View:       c.view,
	}
}

func (c *Pagination) definition() screen.Definition {
	def := c.screen.Definition()
	def.RequireKeys = append(def.RequireKeys, pagination_keys...)
	return def
}

func (c *Pagination) update(state *state.UIState, event screen.ScreenEvent) screen.ScreenResult {
	requiredKey := screen.IsKeyRequired(c.screen.Definition(), event.Key)

	if !requiredKey {
		result := c.localUpdate(state, event)
		if result != nil {
			return *result
		}
	}

	result := c.screen.Update(state, event)
	if !result.IgnoreParents && result.Screen != nil {
		newScreen := NewPagination(*result.Screen).
			ToScreen()
		result.Screen = &newScreen
	}

	return result
}

func (c *Pagination) localUpdate(state *state.UIState, event screen.ScreenEvent) *screen.ScreenResult {
	switch event.Key.Code {
	case key.ActionArrowLeft:
		state.Pager.Page = math.SubClampZero(state.Pager.Page, 1)
		result := screen.ScreenResultFromUIState(state)
		return &result
	case key.ActionArrowRight:
		state.Pager.Page += 1
		result := screen.ScreenResultFromUIState(state)
		return &result
	default:
		return nil
	}
}

func (c *Pagination) view(stt state.UIState) core.ViewModel {
	vm := c.screen.View(stt)

	hasContent := stt.Pager.RestData || stt.Pager.Page > 0
	canShowPage := stt.Pager.ShowPage || vm.IsPagerMode(state.PagerModePage)
	if hasContent && canShowPage {
		page := fmt.Sprintf("page: %d", stt.Pager.Page)

		footer := text.NewLines(
			text.LineJump(),
			text.NewLine(page, style.SpecFromKind(style.SpcKindPaddingRight)),
		)

		vm.Footer.Unshift(
			line.EagerDrawableFromLines(footer...),
		)
	}

	actions := screen.FilterKeyRequired(
		c.screen.Definition(),
		pagination_actions...,
	)

	vm.Helper.Unshift(
		key.ActionsToHelpWithOverride(
			pagination_overrides,
			actions...,
		)...,
	)

	return vm
}
