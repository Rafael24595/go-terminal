package core

import (
	"context"

	assert "github.com/Rafael24595/go-assert/assert/runtime"
	"github.com/Rafael24595/go-terminal/engine/app/cleaner"
	"github.com/Rafael24595/go-terminal/engine/app/screen"
	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/layout"
	"github.com/Rafael24595/go-terminal/engine/model/key"
	"github.com/Rafael24595/go-terminal/engine/model/winsize"
	"github.com/Rafael24595/go-terminal/engine/render"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

type Engine struct {
	running  bool
	context  context.Context
	doneSgnl chan struct{}
	terminal terminal.Terminal
	layout   layout.Layout
	render   render.Render
	cleaner  cleaner.StateCleaner
	screen   screen.Screen
}

func NewEngine(
	terminal terminal.Terminal,
	layout layout.Layout,
	render render.Render,
	cleaner cleaner.StateCleaner,
	screen screen.Screen,
) *Engine {
	return &Engine{
		context:  nil,
		doneSgnl: make(chan struct{}),
		terminal: terminal,
		layout:   layout,
		render:   render,
		cleaner:  cleaner,
		screen:   screen,
	}
}

func (e *Engine) Context(ctx context.Context) *Engine {
	if e.running {
		assert.Unreachable("the engine can be modified after initialization")
		return e
	}

	e.context = ctx

	return e
}

func (e *Engine) Run() <-chan struct{} {
	return e.RunWithContext(
		context.Background(),
	)
}

func (e *Engine) RunWithContext(ctx context.Context) <-chan struct{} {
	if e.running {
		assert.Unreachable("The engine can not be initialized more than once")
		return e.doneSgnl
	}

	e.running = true

	e.context = ctx
	go e.run()

	return e.doneSgnl
}

func (e *Engine) run() {
	defer close(e.doneSgnl)

	e.terminal.OnStart()
	defer e.terminal.OnClose()

	state := state.NewUIState()

	size, err := e.terminal.Size()
	if err != nil {
		println(err)
		return
	}

	e.renderFrame(state, size)

	keys := e.terminal.KeyEvents()
	resizes := e.terminal.ResizeEvents()

	for {
		select {
		//TODO: Add configurable ticker
		case <-e.context.Done():
			return

		case k, ok := <-keys:
			if !ok || k.Code == key.ActionExit {
				return
			}

			e.updateScreen(state, size, k)

		case s, ok := <-resizes:
			if !ok {
				return
			}

			size = s

			e.renderFrame(state, size)

		case <-e.doneSgnl:
			return
		}
	}
}

func (e *Engine) Exit() {
	select {
	case <-e.doneSgnl:
	default:
		close(e.doneSgnl)
	}
}

func (e *Engine) updateScreen(
	state *state.UIState,
	size winsize.Winsize,
	key key.Key,
) *state.UIState {
	result := e.screen.Update(state,
		screen.NewEvent(key),
	)

	state.Pager = result.Pager
	if result.Screen != nil {
		e.screen = *result.Screen
	}

	state = e.cleaner.Cleanup(result, state)
	e.renderFrame(state, size)

	return state
}

func (e *Engine) renderFrame(state *state.UIState, size winsize.Winsize) {
	vmd := e.screen.View(*state)

	lines := e.layout.Apply(state, vmd, size)
	result := e.render.Render(lines, size)

	e.terminal.WriteAll(result)
	e.terminal.Flush()
	e.terminal.Clear()
}
