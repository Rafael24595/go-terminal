package wrapper

import (
	"fmt"

	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/core/drawable/line"
	"github.com/Rafael24595/go-terminal/engine/core/key"
	"github.com/Rafael24595/go-terminal/engine/core/screen"
	"github.com/Rafael24595/go-terminal/engine/core/style"
	"github.com/Rafael24595/go-terminal/engine/core/text"
)

var history_actions = []key.KeyAction{
	key.CustomActionBack,
}

var history_keys = key.NewKeysCode(
	history_actions...,
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
	def.RequireKeys = append(def.RequireKeys, history_keys...)
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

func (c *History) view(state state.UIState) core.ViewModel {
	vm := c.screen.View(state)

	if c.history == nil {
		return vm
	}

	page := fmt.Sprintf("back: %s", c.history.Name())

	footer := text.NewLines(
		text.NewLine(page,
			style.SpecFromKind(style.SpcKindPaddingRight),
		),
	)

	vm.Footer.Unshift(
		line.EagerDrawableFromLines(footer...).
			AddTag(screen.SystemScreenMeta),
	)

	actions := screen.FilterKeyRequired(
		c.screen.Definition(),
		history_actions...,
	)

	vm.Helper.Unshift(
		key.ActionsToHelp(
			actions...,
		)...,
	)

	return vm
}
