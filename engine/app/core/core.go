package core

import (
	"context"
	"time"

	assert "github.com/Rafael24595/go-assert/assert/runtime"
	local "github.com/Rafael24595/go-reacterm-core/engine/commons/log"

	"github.com/Rafael24595/go-log/log"
	"github.com/Rafael24595/go-reacterm-core/engine/app/cleaner"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/layout"
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
	"github.com/Rafael24595/go-reacterm-core/engine/model/pulse"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render"
	"github.com/Rafael24595/go-reacterm-core/engine/terminal"
)

type Engine struct {
	running  bool
	context  context.Context
	doneSgnl chan struct{}
	pulse    *pulse.Pulse
	terminal terminal.Terminal
	layout   layout.Layout
	render   render.Render
	cleaner  cleaner.StateCleaner
	screen   screen.Screen
	passes   []screen.ScreenPass
}

// TODO: Disable pulse on proactive terminal
func NewEngine(
	trm terminal.Terminal,
	lyo layout.Layout,
	rnd render.Render,
	cls cleaner.StateCleaner,
	scn screen.Screen,
) *Engine {
	pulse := pulse.New(50 * time.Millisecond)
	return &Engine{
		context:  nil,
		doneSgnl: make(chan struct{}),
		pulse:    pulse,
		terminal: trm,
		layout:   lyo,
		render:   rnd,
		cleaner:  cls,
		screen:   scn,
		passes:   make([]screen.ScreenPass, 0),
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

func (e *Engine) AddPass(pass ...screen.ScreenPass) *Engine {
	if e.running {
		assert.Unreachable("the engine can be modified after initialization")
		return e
	}

	e.passes = append(e.passes, pass...)
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
	e.compileScreen(e.screen)

	go e.run()

	return e.doneSgnl
}

func (e *Engine) run() {
	defer close(e.doneSgnl)
	defer e.pulse.Exit()

	err := e.terminal.OnStart()
	if err != nil {
		e.logErr(err)
		return
	}

	defer local.LogErrorHandler(e.terminal.OnClose)

	size, err := e.terminal.Size()
	if err != nil {
		e.logErr(err)
		return
	}

	state := state.NewUIState()
	e.renderFrame(state, size)

	keys := e.terminal.KeyEvents()
	resizes := e.terminal.ResizeEvents()

	for {
		select {
		case <-e.context.Done():
			return

		case <-e.pulse.Listen():
			e.renderFrame(state, size)

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

func (e *Engine) compileScreen(screen screen.Screen) *Engine {
	newScreen, err := screen.Compile(e.passes...)
	if err != nil {
		e.logErr(err)
	}

	e.screen = newScreen
	return e
}

func (e *Engine) updateScreen(
	state *state.UIState,
	size winsize.Winsize,
	key key.Key,
) *state.UIState {
	result := e.screen.Update(state,
		screen.NewEvent(key),
	)

	e.manageResult(state, result)
	e.manageScreen(result)

	state = e.cleaner.Cleanup(result, state)

	e.renderFrame(state, size)

	return state
}

func (e *Engine) manageResult(state *state.UIState, result screen.Result) *state.UIState {
	state.Pager = result.Pager
	return state
}

func (e *Engine) manageScreen(result screen.Result) screen.Result {
	if result.Screen != nil {
		e.compileScreen(*result.Screen)
	}
	return result
}

func (e *Engine) syncPulse(vm viewmodel.ViewModel) viewmodel.ViewModel {
	if vm.Behavior.NeedsPulse {
		e.pulse.Enable()
		return vm
	}

	e.pulse.Disable()
	return vm
}

func (e *Engine) renderFrame(state *state.UIState, size winsize.Winsize) {
	vm := e.screen.View(*state)

	lines := e.layout.Apply(state, vm, size)
	result := e.render.Render(lines, size)

	e.syncPager(state, &vm)
	e.syncPulse(vm)

	err := e.terminal.WriteAll(result)
	if err != nil {
		e.logErr(err)
	}

	err = e.terminal.Flush()
	if err != nil {
		e.logErr(err)
	}

	err = e.terminal.Clear()
	if err != nil {
		e.logErr(err)
	}
}

func (e *Engine) syncPager(state *state.UIState, vm *viewmodel.ViewModel) (*state.UIState, *viewmodel.ViewModel) {
	if state.Pager.Syncronyzed {
		return state, vm
	}

	vm.Behavior.NeedsPulse = true
	state.Pager.Syncronyzed = true
	return state, vm
}

func (e *Engine) logErr(err error) {
	log.Error(err)
	assert.Unreachable("error: %s", err)
}
