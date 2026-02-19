package commons

import (
	"fmt"

	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/core/drawable/line"
	"github.com/Rafael24595/go-terminal/engine/core/screen"
	"github.com/Rafael24595/go-terminal/engine/core/style"
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
		Definition: c.screen.Definition,
		Update:     c.update,
		View:       c.view,
	}
}

func (c *History) update(state state.UIState, event screen.ScreenEvent) screen.ScreenResult {
	requiredKey := isKeyRequired(c.screen.Definition(), event.Key)

	if !requiredKey {
		result := c.localUpdate(state, event)
		if result != nil {
			return *result
		}
	}

	result := c.screen.Update(state, event)
	if result.Screen != nil {
		newBack := NewHistory(*result.Screen)
		newBack.history = &c.screen
		newScreen := newBack.ToScreen()
		result.Screen = &newScreen
	}
	return result
}

func (c *History) localUpdate(_ state.UIState, event screen.ScreenEvent) *screen.ScreenResult {
	if event.Key.Rune == 'b' && c.history != nil {
		newBack := NewHistory(*c.history)
		newScreen := newBack.ToScreen()
		result := screen.ScreenResultFromScreen(&newScreen)
		return &result
	}

	return nil
}

func (c *History) view(state state.UIState) core.ViewModel {
	vm := c.screen.View(state)

	if c.history == nil {
		return vm
	}

	page := fmt.Sprintf("back: %s", c.history.Name())

	footer := core.NewLines(
		core.LineJump(),
		core.NewLine(page, style.SpecFromKind(style.SpcKindPaddingRight)),
	)

	vm.Footer.Unshift(
		line.LinesEagerDrawableFromLines(footer...),
	)

	return vm
}
