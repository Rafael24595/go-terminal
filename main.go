package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	assert "github.com/Rafael24595/go-assert/assert/runtime"
	
	"github.com/Rafael24595/go-log/log"
	"github.com/Rafael24595/go-log/log/model/record"
	"github.com/Rafael24595/go-log/log/provider/file"
	"github.com/Rafael24595/go-terminal/engine/app/runtime"
	"github.com/Rafael24595/go-terminal/engine/app/screen"
	"github.com/Rafael24595/go-terminal/engine/app/screen/partial"
	"github.com/Rafael24595/go-terminal/engine/app/screen/wrapper"
	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/layout"
	"github.com/Rafael24595/go-terminal/engine/model/action"
	"github.com/Rafael24595/go-terminal/engine/model/inline"
	"github.com/Rafael24595/go-terminal/engine/model/key"
	"github.com/Rafael24595/go-terminal/engine/render"
	"github.com/Rafael24595/go-terminal/engine/render/spacer"
	"github.com/Rafael24595/go-terminal/engine/terminal"

	context_cleaner "github.com/Rafael24595/go-terminal/engine/app/cleaner/context"

	wrapper_layout "github.com/Rafael24595/go-terminal/wrapper/layout"
	wrapper_render "github.com/Rafael24595/go-terminal/wrapper/render"

	wrapper_screen "github.com/Rafael24595/go-terminal/wrapper/screen"
	wrapper_terminal "github.com/Rafael24595/go-terminal/wrapper/terminal"
)

const paddingCols = 10
const paddingRows = 5

// Move main and wrapper packages to a new project
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	provider := file.FileProvider{
		Session: runtime.Instance.SessionId(),
	}

	if err := log.DefaultFromProvider(ctx, provider); err != nil {
		panic(err.Error())
	}

	defer log.OnClose()

	assert.DefaultWriter(
		log.WriterFromCategory(record.WARNING),
	)

	state := &state.UIState{}

	cmd := wrapper_terminal.NewConsole()
	cmd.Color("\x1b[0;32m")

	trm := cmd.ToTerminal()

	size, _ := trm.Size()

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

	lyt := layout.NewLayout(wrapper_layout.TerminalApply)
	lytf := wrapper_layout.NewFixed(lyt, pr, pc)
	lyt = lytf.ToLayout()

	rnd := render.NewRender(wrapper_render.TerminalRender)
	rndf := wrapper_render.NewFixed(rnd, pr, pc)
	rnd = rndf.ToRender()

	cls := context_cleaner.NewContextCleaner()

	trm.OnStart()
	defer trm.OnClose()

	inputChan := make(chan key.Key, 64)
	go readInput(trm, inputChan)

	//TODO: Use events instead sleep loop
	for {
		newSize, _ := trm.Size()

		//TODO: Replace with chan events
		if newSize.Cols != size.Cols || newSize.Rows != size.Rows {
			cmd.Update()
			lytf.Update(newSize.Rows-paddingRows, newSize.Cols-paddingCols)
			rndf.Update(newSize.Rows-paddingRows, newSize.Cols-paddingCols)
		}

		size = newSize

		vmd := stc.View(*state)

		lns := lyt.Apply(state, vmd, size)
		str := rnd.Render(lns, size)

		trm.WriteAll(str)
		trm.Flush()
		trm.Clear()

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

			state = cls.Cleanup(result, state)

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
