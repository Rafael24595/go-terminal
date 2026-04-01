package wrapper

import (
	"fmt"

	assert "github.com/Rafael24595/go-assert/assert/runtime"
	"github.com/Rafael24595/go-terminal/engine/app/pager"
	"github.com/Rafael24595/go-terminal/engine/app/screen"
	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/app/viewmodel"
	"github.com/Rafael24595/go-terminal/engine/helper/math"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable/line"
	"github.com/Rafael24595/go-terminal/engine/model/help"
	"github.com/Rafael24595/go-terminal/engine/model/key"
	"github.com/Rafael24595/go-terminal/engine/render/style"
	"github.com/Rafael24595/go-terminal/engine/render/text"
)

var sources = map[pager.EngineCode]screen.DefinitionSources{
	pager.CodeEnginePaged:  pager_definition,
	pager.CodeEngineScroll: scroll_definition,
}

var pager_definition = screen.NewDefinitionSources(
	map[key.KeyAction]help.HelpField{
		key.ActionArrowLeft:  {Code: []string{"←"}, Detail: "Prev page"},
		key.ActionArrowRight: {Code: []string{"→"}, Detail: "Next page"},
	},
	[]key.KeyAction{
		key.ActionArrowLeft,
		key.ActionArrowRight,
	},
)

var scroll_definition = screen.NewDefinitionSources(
	map[key.KeyAction]help.HelpField{
		key.ActionArrowUp:   {Code: []string{"↑"}, Detail: "Scroll up"},
		key.ActionArrowDown: {Code: []string{"↓"}, Detail: "Scroll down"},
	},
	[]key.KeyAction{
		key.ActionArrowUp,
		key.ActionArrowDown,
	},
)

type Pagination struct {
	engine      pager.EngineCode
	screen      screen.Screen
	forceEngine *pager.Engine
}

func NewPagination(screen screen.Screen) *Pagination {
	return &Pagination{
		engine:      pager.CodeEnginePaged,
		screen:      screen,
		forceEngine: nil,
	}
}

func (c *Pagination) ForceEngine(forceEngine pager.Engine) *Pagination {
	c.forceEngine = &forceEngine
	c.engine = forceEngine.Code
	return c
}

func (c *Pagination) ToScreen() screen.Screen {
	return screen.Screen{
		Name:       c.screen.Name,
		Definition: c.definition,
		Update:     c.update,
		View:       c.view,
		Stack:      c.screen.Stack,
	}
}

func (c *Pagination) definitionSource() screen.DefinitionSources {
	source, ok := sources[c.engine]
	if ok {
		return source
	}

	assert.Unreachable("unhandled engine definition %d", c.engine)

	return pager_definition
}

func (c *Pagination) definition() screen.Definition {
	def := c.screen.Definition()
	def.RequireKeys = append(def.RequireKeys, c.definitionSource().Keys...)
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
	if result.Screen != nil {
		newWrapper := NewPagination(*result.Screen)
		newWrapper.engine = c.engine
		newWrapper.forceEngine = c.forceEngine
		newScreen := newWrapper.ToScreen()
		result.Screen = &newScreen
	}

	return result
}

func (c *Pagination) localUpdate(state *state.UIState, event screen.ScreenEvent) *screen.ScreenResult {
	keyback := key.ActionArrowLeft
	keyNext := key.ActionArrowRight
	if c.engine == pager.CodeEngineScroll {
		keyback = key.ActionArrowUp
		keyNext = key.ActionArrowDown
	}

	if event.Key.Code == keyback {
		state.Pager.Page = math.SubClampZero(state.Pager.Page, 1)
		result := screen.ScreenResultFromUIState(state)

		return &result
	}

	if event.Key.Code == keyNext {
		state.Pager.Page += 1
		result := screen.ScreenResultFromUIState(state)

		return &result
	}

	return nil
}

func (c *Pagination) view(stt state.UIState) viewmodel.ViewModel {
	vm := c.screen.View(stt)

	if c.forceEngine != nil {
		vm.Pager.SetEngine(*c.forceEngine)
	}

	source := c.definitionSource()

	if c.shouldShowPage(stt, vm) {
		label := "page"
		if vm.Pager.Engine.Code == pager.CodeEngineScroll {
			label = "scroll"
		}

		footer := text.NewLines(
			text.NewLine(
				fmt.Sprintf("%s: %d", label, stt.Pager.Page),
				style.SpecFromKind(style.SpcKindPaddingRight),
			),
		)

		vm.Footer.Unshift(
			line.EagerDrawableFromLines(footer...).
				AddTag(screen.SystemScreenMeta),
		)
	}

	actions := screen.FilterKeyRequired(
		c.screen.Definition(),
		source.Actions...,
	)

	vm.Helper.Unshift(
		key.ActionsToHelpWithOverride(
			source.Overrides,
			actions...,
		)...,
	)

	return vm
}

func (c *Pagination) shouldShowPage(stt state.UIState, vm viewmodel.ViewModel) bool {
	predicate := vm.Pager.Predicate.Code
	//engine := vm.Pager.Engine.Code

	if predicate != pager.CodePredicatePage /*|| engine == pager.CodeEngineScroll*/ {
		return false
	}

	if stt.Pager.ForceShow {
		return true
	}

	return stt.Pager.HasMore || stt.Pager.Page > 0
}
