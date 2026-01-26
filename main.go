package main

import (
	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/render"
	wrapper_layout "github.com/Rafael24595/go-terminal/wrapper/layout"
	wrapper_render "github.com/Rafael24595/go-terminal/wrapper/render"
	wrapper_screen "github.com/Rafael24595/go-terminal/wrapper/screen"
	wrapper_console "github.com/Rafael24595/go-terminal/wrapper/terminal"
)

// Move main and wrapper packages to a new project
func main() {
	state := state.UIState{
		Page: 0,
	}

	t := wrapper_console.NewConsole().ToTerminal()

	i := wrapper_screen.NewIndex().ToScreen()
	s := wrapper_screen.NewWrapperMain(i).ToScreen()

	l := core.NewLayout(wrapper_layout.TerminalApply)
	l = wrapper_layout.NewFixed(l, 40, 100).ToLayout()

	r := render.NewRender(wrapper_render.TerminalRender)
	r = wrapper_render.NewFixed(r, 40, 100).ToRender()

	t.OnStart()
	defer t.OnClose()

	for {
		vmd := s.View()

		lns := l.Apply(&state, vmd, t.Size())
		str := r.Render(lns, t.Size())

		t.WriteAll(str)

		t.Flush()

		t.Clear()
	}
}
