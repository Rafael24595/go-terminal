package main

import (
	"time"

	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/core/drawable/action"
	"github.com/Rafael24595/go-terminal/engine/core/inline"
	"github.com/Rafael24595/go-terminal/engine/core/key"
	"github.com/Rafael24595/go-terminal/engine/core/screen"
	"github.com/Rafael24595/go-terminal/engine/core/screen/partial"
	"github.com/Rafael24595/go-terminal/engine/core/screen/wrapper"
	"github.com/Rafael24595/go-terminal/engine/core/spacer"
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
	state := &state.UIState{}

	c := wrapper_terminal.NewConsole()
	c.Color("\x1b[0;32m")

	t := c.ToTerminal()

	size, _ := t.Size()

	pc := size.Cols - paddingCols
	pr := size.Rows - paddingRows

	lnd := wrapper_screen.NewLanding()
	hdr := wrapper_screen.NewBaseHeader(lnd)

	his := wrapper.NewHistory(hdr).ToScreen()
	pge := wrapper.NewPagination(his).ToScreen()
	hlp := wrapper.NewHelp(pge).ToScreen()

	inl := partial.NewInline(hlp).
		PushAction(action.FocusFooter,
			inline.NewFilterMeta(inline.TargetTags, screen.SystemScreenMeta),
		).
		ToScreen()

	stc := partial.NewSpacer(inl).
		Header(spacer.NewSpacerMeta(1, spacer.SpacerAfterEach)).
		Footer(spacer.NewSpacerMeta(1, spacer.SpacerAfterEach)).
		ToScreen()

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
		newSize, _ := t.Size()

		//TODO: Replace with chan events
		if newSize.Cols != size.Cols || newSize.Rows != size.Rows {
			c.Update()
			lf.Update(newSize.Rows-paddingRows, newSize.Cols-paddingCols)
			rf.Update(newSize.Rows-paddingRows, newSize.Cols-paddingCols)
		}

		size = newSize

		vmd := stc.View(*state)

		lns := l.Apply(state, vmd, size)
		str := r.Render(lns, size)

		t.WriteAll(str)

		t.Flush()

		t.Clear()

		select {
		case key, ok := <-inputChan:
			if !ok {
				return
			}
			result := stc.Update(state, screen.ScreenEvent{
				Key: key,
			})

			state.Pager = result.Pager

			if result.Screen != nil {
				stc = *result.Screen
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
		if rn.Code == key.ActionExit {
			close(ch)
			return
		}
	}
}
