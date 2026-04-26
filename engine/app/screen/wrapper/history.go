package wrapper

import (
	"fmt"

	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/block"
	"github.com/Rafael24595/go-reacterm-core/engine/model/help"
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

var history_definition = screen.NewDefinitionSources(
	map[key.KeyAction]help.HelpField{},
	[]key.KeyAction{
		key.CustomActionBack,
	},
)

type History struct {
	history *screen.Screen
	screen  screen.Screen
}

func NewHistory(screen screen.Screen) *History {
	return &History{
		screen: screen,
	}
}

func (c *History) ToScreen() screen.Screen {
	return screen.Screen{
		Name:       c.screen.Name,
		Definition: c.definition,
		Update:     c.update,
		View:       c.view,
		Stack:      c.screen.Stack,
	}
}

func (c *History) definition() screen.Definition {
	def := c.screen.Definition()
	def.RequireKeys = append(def.RequireKeys, history_definition.Keys...)
	return def
}

func (c *History) update(state *state.UIState, event screen.ScreenEvent) screen.ScreenResult {
	requiredKey := screen.IsKeyRequired(c.screen.Definition(), event.Key)

	if !requiredKey {
		result := c.localUpdate(state, event)
		if result != nil {
			return *result
		}
	}

	result := c.screen.Update(state, event)
	if result.Screen != nil {
		newWrapper := NewHistory(*result.Screen)
		newWrapper.history = &c.screen
		newScreen := newWrapper.ToScreen()
		result.Screen = &newScreen
	}

	return result
}

func (c *History) localUpdate(_ *state.UIState, event screen.ScreenEvent) *screen.ScreenResult {
	if c.history == nil || event.Key.Code != key.CustomActionBack {
		return nil
	}

	newBack := NewHistory(*c.history)
	newScreen := newBack.ToScreen()
	result := screen.ScreenResultFromScreen(&newScreen)

	return &result
}

func (c *History) view(state state.UIState) viewmodel.ViewModel {
	vm := c.screen.View(state)

	if c.history == nil {
		return vm
	}

	page := fmt.Sprintf("back: %s", c.history.Name())

	footer := []text.Line{
		*text.NewLine(page,
			style.SpecFromKind(style.SpcKindPaddingRight),
		),
	}

	vm.Footer.Unshift(
		block.BlockDrawableFromLines(footer...).
			AddTag(screen.SystemMetaTag),
	)

	actions := screen.FilterKeyRequired(
		c.screen.Definition(),
		history_definition.Actions...,
	)

	vm.Helper.Unshift(
		key.ActionsToHelp(
			actions...,
		)...,
	)

	return vm
}
