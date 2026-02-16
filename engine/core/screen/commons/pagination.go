package commons

import (
	"fmt"

	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/core/drawable/line"
	"github.com/Rafael24595/go-terminal/engine/core/key"
	"github.com/Rafael24595/go-terminal/engine/core/screen"
	"github.com/Rafael24595/go-terminal/engine/helper/math"
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
		Definition: c.screen.Definition,
		Update:     c.Update,
		View:       c.View,
	}
}

func (c *Pagination) Update(state state.UIState, event screen.ScreenEvent) screen.ScreenResult {
	requiredKey := isKeyRequired(c.screen.Definition(), event.Key)

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

func (c *Pagination) localUpdate(state state.UIState, event screen.ScreenEvent) *screen.ScreenResult {
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

func (c *Pagination) View(state state.UIState) core.ViewModel {
	vm := c.screen.View(state)

	if vm.Pager.Enabled {
		page := fmt.Sprintf("page: %d", vm.Pager.Page)

		footer := core.NewLines(
			core.LineJump(),
			core.NewLine(page, core.ModePadding(core.Right)),
		)

		vm.Footer.Unshift(
			line.LinesEagerDrawableFromLines(footer...),
		)
	}
	
	return vm
}
