package main

import (
	"time"

	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/core/key"
	"github.com/Rafael24595/go-terminal/engine/core/screen"
	"github.com/Rafael24595/go-terminal/engine/core/screen/commons"
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
		Pager: state.PagerState{
			Page:    0,
			Enabled: false,
		},
		Cursor: state.CursorState{
			Cursor: 0,
			Offset: 0,
		},
	}

	c := wrapper_terminal.NewConsole()
	c.Color("\x1b[0;32m")

	t := c.ToTerminal()

	size := t.Size()

	pc := size.Cols - paddingCols
	pr := size.Rows - paddingRows

	i := wrapper_screen.NewLanding()
	p := commons.NewPagination(i).ToScreen()
	h := commons.NewHistory(p).ToScreen()
	s := wrapper_screen.NewBaseHeader(h)

	l := core.NewLayout(wrapper_layout.TerminalApply)
	lf := wrapper_layout.NewFixed(l, pr, pc)
	l = lf.ToLayout()

	r := render.NewRender(wrapper_render.TerminalRender)
	rf := wrapper_render.NewFixed(r, pr, pc)
	r = rf.ToRender()

	t.OnStart()
	defer t.OnClose()

	inputChan := make(chan key.Key, 64)
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

		vmd := s.View(state)

		state.Pager = vmd.Pager
		state.Cursor = vmd.Cursor

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
			result := s.Update(state, screen.ScreenEvent{
				Key: key,
			})

			state.Pager = result.Pager
			state.Cursor = result.Cursor

			if result.Screen != nil {
				s = *result.Screen
			}

		default:
		}

		time.Sleep(20 * time.Millisecond)
	}
}

func readInput(t terminal.Terminal, ch chan<- key.Key) {
	for {
		rn, err := t.ReadKey()
		if err != nil {
			println(err)
			close(ch)
			return
		}

		ch <- *rn
		if rn.Code == key.KeyCtrlC {
			close(ch)
			return
		}
	}
}
