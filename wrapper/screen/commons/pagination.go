package wrapper_commons

import (
	"fmt"

	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/helper/math"
	wrapper_terminal "github.com/Rafael24595/go-terminal/wrapper/terminal"
)

type Pagination struct {
	screen core.Screen
}

func NewPagination(screen core.Screen) *Pagination {
	return &Pagination{
		screen: screen,
	}
}

func (c *Pagination) ToScreen() core.Screen {
	return core.Screen{
		Name:   c.screen.Name,
		Update: c.Update,
		View:   c.View,
	}
}

func (c *Pagination) Update(state state.UIState, event core.ScreenEvent) core.ScreenResult {
	if result := c.localUpdate(state, event); result != nil {
		return *result
	}

	result := c.screen.Update(state, event)
	if !result.IgnoreParents && result.Screen != nil {
		newScreen := NewPagination(*result.Screen).
			ToScreen()
		result.Screen = &newScreen
	}

	return result
}

func (c *Pagination) localUpdate(state state.UIState, event core.ScreenEvent) *core.ScreenResult {
	switch event.Key {
	case wrapper_terminal.ARROW_LEFT:
		state.Layout.Page = math.SubClampZero(state.Layout.Page, 1)
		result := core.ScreenResultFromState(state)
		return &result
	case wrapper_terminal.ARROW_RIGHT:
		state.Layout.Page += 1

		result := core.ScreenResultFromState(state)
		return &result
	default:
		return nil
	}
}

func (c *Pagination) View(state state.UIState) core.ViewModel {
	vm := c.screen.View(state)
	if state.Layout.Pagination {
		page := fmt.Sprintf("page: %d", state.Layout.Page)
		footer := core.NewLines(
			core.LineJump(),
			core.NewLine(page, core.ModePadding(core.Right)),
		)
		vm.Footer = append(vm.Footer, footer...)
	}
	return vm
}
