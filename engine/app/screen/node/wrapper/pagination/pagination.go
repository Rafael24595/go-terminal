package pagination

import (
	"fmt"

	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/app/pager"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/drain"
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

const errf_unhandled = "unhandled pager type '%d'"

var base_definition = screen.NewDefinition(
	map[key.Action]key.Descriptor{
		key.ActionPageUp:   {Code: []string{"⇞"}, Detail: "Prev page"},
		key.ActionPageDown: {Code: []string{"⇟"}, Detail: "Next page"},
	},
	[]key.Action{
		key.ActionPageUp,
		key.ActionPageDown,
	},
)

var definitions = map[pager.EngineCode]screen.Definition{
	pager.CodeEnginePaged:  pager_definition,
	pager.CodeEngineScroll: scroll_definition,
}

var keys = map[pager.EngineCode]struct {
	back key.Action
	next key.Action
}{
	pager.CodeEnginePaged:  {key.ActionArrowLeft, key.ActionArrowRight},
	pager.CodeEngineScroll: {key.ActionArrowUp, key.ActionArrowDown},
}

var labels = map[pager.EngineCode]string{
	pager.CodeEnginePaged:  "page",
	pager.CodeEngineScroll: "scroll",
}

var pager_definition = screen.NewDefinition(
	map[key.Action]key.Descriptor{
		key.ActionArrowLeft:  {Code: []string{"←"}, Detail: "Prev page"},
		key.ActionArrowRight: {Code: []string{"→"}, Detail: "Next page"},
	},
	[]key.Action{
		key.ActionArrowLeft,
		key.ActionArrowRight,
	},
)

var scroll_definition = screen.NewDefinition(
	map[key.Action]key.Descriptor{
		key.ActionArrowUp:   {Code: []string{"↑"}, Detail: "Scroll up"},
		key.ActionArrowDown: {Code: []string{"↓"}, Detail: "Scroll down"},
	},
	[]key.Action{
		key.ActionArrowUp,
		key.ActionArrowDown,
	},
)

type Pagination struct {
	engine      pager.EngineCode
	node        screen.Node
	forceEngine *pager.Engine
}

func New(screen screen.Node) *Pagination {
	return &Pagination{
		engine:      pager.CodeEnginePaged,
		node:        screen,
		forceEngine: nil,
	}
}

func (c *Pagination) ForceEngine(forceEngine pager.Engine) *Pagination {
	c.forceEngine = &forceEngine
	c.engine = forceEngine.Code
	return c
}

func (c *Pagination) ToNode() screen.Node {
	return screen.NewBuilder().
		Name(c.node.Screen.Name).
		AddStack(c.node.Stack).
		Definition(c.definition).
		Update(c.update).
		View(c.view).
		Children(c.node).
		ToNode()
}

func (c *Pagination) definition() screen.Definition {
	node := c.node.Screen.Definition()
	return base_definition.Merge(
		c.findDefinition().Merge(node),
	)
}

func (c *Pagination) findDefinition() screen.Definition {
	if source, ok := definitions[c.engine]; ok {
		return source
	}

	assert.Unreachable("unhandled engine definition %d", c.engine)
	return pager_definition
}

func (c *Pagination) update(state *state.UIState, event screen.Event) screen.Result {
	definition := c.node.Screen.Definition()

	if !definition.IsRequired(event.Key) {
		result := c.localUpdate(state, event)
		if result != nil {
			return *result
		}
	}

	result := c.node.Screen.Update(state, event)
	if result.Node == nil {
		return result
	}

	newWrapper := New(*result.Node)
	newWrapper.engine = c.engine
	newWrapper.forceEngine = c.forceEngine
	newNode := newWrapper.ToNode()
	result.Node = &newNode

	return result
}

func (c *Pagination) localUpdate(state *state.UIState, event screen.Event) *screen.Result {
	keys, ok := keys[pager.CodeEnginePaged]

	assert.True(ok, errf_unhandled, pager.CodeEnginePaged)

	if event.Key.Code == key.ActionPageUp || event.Key.Code == keys.back {
		state.Pager.DecTarget()
		result := screen.ResultFromUIState(state)
		return &result
	}

	if event.Key.Code == key.ActionPageDown || event.Key.Code == keys.next {
		state.Pager.IncTarget()
		result := screen.ResultFromUIState(state)
		return &result
	}

	return nil
}

func (c *Pagination) view(stt state.UIState) viewmodel.ViewModel {
	vm := c.node.Screen.View(stt)
	if c.forceEngine != nil {
		vm.Pager.SetEngine(*c.forceEngine)
	}

	if !c.shouldShowPage(stt, vm) {
		return vm
	}

	label, ok := labels[pager.CodeEnginePaged]

	assert.True(ok, errf_unhandled, pager.CodeEnginePaged)

	footer := []text.Line{
		*text.NewLine(
			fmt.Sprintf("%s: %d", label, stt.Pager.ActualPage),
			style.SpecFromKind(style.SpcKindPaddingRight),
		),
	}

	vm.Footer.Unshift(
		drain.DrawableFromLines(footer...).
			AddTag(screen.SystemMetaTag),
	)

	return vm
}

func (c *Pagination) shouldShowPage(stt state.UIState, vm viewmodel.ViewModel) bool {
	predicate := vm.Pager.Predicate.Code

	if predicate != pager.CodePredicatePage {
		return false
	}

	if stt.Pager.ForceShow {
		return true
	}

	return stt.Pager.HasMore || stt.Pager.ActualPage > 0
}
