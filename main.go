package main

import (
	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/render"
	"github.com/Rafael24595/go-terminal/engine/terminal"
	wrapper_layout "github.com/Rafael24595/go-terminal/wrapper/layout"
	wrapper_render "github.com/Rafael24595/go-terminal/wrapper/render"

	wrapper_screen "github.com/Rafael24595/go-terminal/wrapper/screen"
	wrapper_terminal "github.com/Rafael24595/go-terminal/wrapper/terminal"
)

const paddingCols = 10
const paddingRows = 5

// Move main and wrapper packages to a new project
func main() {
	state := state.UIState{
		Page: 0,
	}

	c := wrapper_terminal.NewConsole()
	c.Color("\x1b[0;32m")

	t := c.ToTerminal()

	size := t.Size()

	pc := size.Cols - paddingCols
	pr := size.Rows - paddingRows

	i := wrapper_screen.NewLanding()
	s := wrapper_screen.NewWrapperMain(i).ToScreen()

	l := core.NewLayout(wrapper_layout.TerminalApply)
	lf := wrapper_layout.NewFixed(l, pr, pc)
	l = lf.ToLayout()

	r := render.NewRender(wrapper_render.TerminalRender)
	rf := wrapper_render.NewFixed(r, pr, pc)
	r = rf.ToRender()

	t.OnStart()
	defer t.OnClose()

	inputChan := make(chan string)
	go readInput(t, inputChan)

	for {
		newSize := t.Size()

		//TODO: Replace with chan events
		if newSize.Cols != size.Cols || newSize.Rows != size.Rows {
			c.Update()
			lf.Update(newSize.Rows-paddingRows, newSize.Cols-paddingCols)
			rf.Update(newSize.Rows-paddingRows, newSize.Cols-paddingCols)
		}

		size = newSize

		vmd := s.View()

		lns := l.Apply(&state, vmd, size)
		str := r.Render(lns, size)

		t.WriteAll(str)

		t.Flush()

		t.Clear()

		select {
		case key, ok := <-inputChan:
			if !ok {
				return
			}
			s.Update(core.ScreenEvent{
				Key: key,
			})
		default:
		}
	}
}

func readInput(t terminal.Terminal, ch chan<- string) {
	for {
		rn, err := t.ReadKey()
		if err != nil {
			println(err)
			close(ch)
			return
		}

		ch <- rn
		if rn == wrapper_terminal.CTRL_C {
			close(ch)
			return
		}
	}
}
