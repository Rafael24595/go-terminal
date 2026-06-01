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

func (n *Pagination) ForceEngine(forceEngine pager.Engine) *Pagination {
	n.forceEngine = &forceEngine
	n.engine = forceEngine.Code
	return n
}

func (n *Pagination) ToNode() screen.Node {
	return screen.NewBuilder().
		Name(n.node.Name).
		AddStack(n.node.Stack).
		Keys(n.keys).
		Tick(n.tick).
		View(n.view).
		Children(n.node).
		ToNode()
}

func (n *Pagination) keys() screen.Definition {
	node := n.node.Screen.Keys()
	return base_definition.Merge(
		n.findDefinition().Merge(node),
	)
}

func (n *Pagination) findDefinition() screen.Definition {
	if source, ok := definitions[n.engine]; ok {
		return source
	}

	assert.Unreachable("unhandled engine definition %d", n.engine)
	return pager_definition
}

func (n *Pagination) tick(uiState *state.UIState, event screen.Event) screen.Result {
	definition := n.node.Screen.Keys()

	if !definition.IsRequired(event.Key) {
		result := n.localTick(uiState, event)
		if result != nil {
			return *result
		}
	}

	result := n.node.Screen.Tick(uiState, event)
	if result.Node == nil {
		return result
	}

	newWrapper := New(*result.Node)
	newWrapper.engine = n.engine
	newWrapper.forceEngine = n.forceEngine
	newNode := newWrapper.ToNode()
	result.Node = &newNode

	return result
}

func (n *Pagination) localTick(uiState *state.UIState, event screen.Event) *screen.Result {
	keys, ok := keys[pager.CodeEnginePaged]

	assert.True(ok, errf_unhandled, pager.CodeEnginePaged)

	if event.Key.Code == key.ActionPageUp || event.Key.Code == keys.back {
		uiState.Pager.DecTarget()
		result := screen.ResultFromUIState(uiState)
		return &result
	}

	if event.Key.Code == key.ActionPageDown || event.Key.Code == keys.next {
		uiState.Pager.IncTarget()
		result := screen.ResultFromUIState(uiState)
		return &result
	}

	return nil
}

func (n *Pagination) view(uiState state.UIState) viewmodel.ViewModel {
	vm := n.node.Screen.View(uiState)
	if n.forceEngine != nil {
		vm.Pager.SetEngine(*n.forceEngine)
	}

	if !n.shouldShowPage(uiState, vm) {
		return vm
	}

	label, ok := labels[vm.Pager.Engine.Code]

	assert.True(ok, errf_unhandled, vm.Pager.Engine.Code)

	footer := []text.Line{
		*text.NewLine(
			fmt.Sprintf("%s: %d", label, uiState.Pager.ActualPage),
			style.SpecFromKind(style.SpcKindPaddingRight),
		),
	}

	vm.Footer.Unshift(
		drain.UnitFromLines(footer...).
			AddTag(screen.SystemMetaTag),
	)

	return vm
}

func (n *Pagination) shouldShowPage(uiState state.UIState, vm viewmodel.ViewModel) bool {
	predicate := vm.Pager.Predicate.Code

	if predicate != pager.CodePredicatePage {
		return false
	}

	if uiState.Pager.ForceShow {
		return true
	}

	return uiState.Pager.HasMore || uiState.Pager.ActualPage > 0
}
