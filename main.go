package main

import (
	"github.com/Rafael24595/go-terminal/engine/render"
	wrapper_render "github.com/Rafael24595/go-terminal/wrapper/render"
	wrapper_screen "github.com/Rafael24595/go-terminal/wrapper/screen"
	wrapper_console "github.com/Rafael24595/go-terminal/wrapper/terminal"
)

//Move main and wrapper packages to a new project
func main() {
	t := wrapper_console.NewConsole().ToTerminal()
	s := wrapper_screen.NewWrapperMain().ToScreen()
	r := render.Render{
		Render: wrapper_render.TerminalRender,
	}

	t.OnStart()
	defer t.OnClose()

	for {
		vm := s.View()

		l := r.Render(vm, t.Size())

		t.WriteAll(l)

		t.Flush()

		t.Clear()
	}
}
