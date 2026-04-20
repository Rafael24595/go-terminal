package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-log/log"
	"github.com/Rafael24595/go-log/log/model/record"
	"github.com/Rafael24595/go-log/log/provider/file"
	"github.com/Rafael24595/go-terminal/engine/app/core"
	"github.com/Rafael24595/go-terminal/engine/app/pager"
	"github.com/Rafael24595/go-terminal/engine/app/runtime"
	"github.com/Rafael24595/go-terminal/engine/app/screen"
	"github.com/Rafael24595/go-terminal/engine/app/screen/partial"
	"github.com/Rafael24595/go-terminal/engine/app/screen/wrapper"
	"github.com/Rafael24595/go-terminal/engine/layout"
	"github.com/Rafael24595/go-terminal/engine/model/action"
	"github.com/Rafael24595/go-terminal/engine/model/inline"
	"github.com/Rafael24595/go-terminal/engine/model/winsize"
	"github.com/Rafael24595/go-terminal/engine/model/winsize/transformer"
	"github.com/Rafael24595/go-terminal/engine/render"
	"github.com/Rafael24595/go-terminal/engine/render/adapter"
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

	configLog(ctx)
	defer log.OnClose()

	terminal := makeTerminal(ctx)

	transformer := transformer.WithMargin(paddingRows, paddingCols)
	layout := makeLayout(transformer)
	render := makeRender(transformer)

	cleaner := context_cleaner.NewContextCleaner()

	screen := makeScreen()

	<-core.NewEngine(
		terminal,
		layout,
		render,
		cleaner,
		screen,
	).RunWithContext(ctx)
}

func configLog(ctx context.Context) {
	provider := file.FileProvider{
		Session: runtime.Instance.SessionId(),
	}

	if err := log.DefaultFromProvider(ctx, provider); err != nil {
		panic(err.Error())
	}

	assert.DefaultWriter(
		log.WriterFromCategory(record.WARNING),
	)
}

func makeTerminal(ctx context.Context) terminal.Terminal {
	return wrapper_terminal.NewConsole().
		Color("\x1b[0;32m").
		Context(ctx).
		ToTerminal()
}

func makeScreen() screen.Screen {
	landing := wrapper_screen.NewLanding()
	header := wrapper_screen.NewBaseHeader(landing)

	history := wrapper.NewHistory(header).ToScreen()
	pagination := wrapper.NewPagination(history).
		ForceEngine(pager.EnginePage()).
		ToScreen()
	helper := wrapper.NewHelp(pagination).ToScreen()

	inline := partial.NewInline(helper).
		PushAction(action.FocusFooter,
			inline.NewFilterMeta(inline.TargetTags, screen.SystemScreenMeta),
		).
		ToScreen()

	return partial.NewSpacer(inline).
		Header(spacer.NewSpacerMeta(1, spacer.SpacerAfterEach)).
		Footer(spacer.NewSpacerMeta(1, spacer.SpacerAfterEach)).
		ToScreen()
}

func makeLayout(transformer winsize.Transformer) layout.Layout {
	return layout.NewBuilder(wrapper_layout.TerminalApply).
		Transformer(transformer).
		ToLayout()
}

func makeRender(transformer winsize.Transformer) render.Render {
	adapter := adapter.WithPadding(
		transformer,
		wrapper_render.TerminalRawRender,
	)

	return render.NewBuilder(adapter).
		ToRender()
}
