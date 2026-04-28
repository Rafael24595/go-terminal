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
	"github.com/Rafael24595/go-reacterm-core/engine/app/core"
	"github.com/Rafael24595/go-reacterm-core/engine/app/pager"
	"github.com/Rafael24595/go-reacterm-core/engine/app/runtime"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/partial/pipeline"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/partial/pipeline/inline"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/partial/pipeline/spacer"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/wrapper/help"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/wrapper/history"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/wrapper/pagination"
	"github.com/Rafael24595/go-reacterm-core/engine/layout"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize/transformer"
	"github.com/Rafael24595/go-reacterm-core/engine/render"
	"github.com/Rafael24595/go-reacterm-core/engine/render/adapter"
	"github.com/Rafael24595/go-reacterm-core/engine/terminal"

	"github.com/Rafael24595/go-reacterm-core/engine/app/cleaner/composite"
	"github.com/Rafael24595/go-reacterm-core/engine/app/cleaner/stack"

	local "github.com/Rafael24595/go-reacterm-core/engine/commons/log"

	wrapper_layout "github.com/Rafael24595/go-reacterm-core/wrapper/layout"
	wrapper_render "github.com/Rafael24595/go-reacterm-core/wrapper/render"
	wrapper_console "github.com/Rafael24595/go-reacterm-core/wrapper/terminal/console"

	wrapper_screen "github.com/Rafael24595/go-reacterm-core/wrapper/screen"
)

const paddingCols = 10
const paddingRows = 5

// Move main and wrapper packages to a new project
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	configLog(ctx)
	defer local.WriterErrorHandler(os.Stderr, log.OnClose)

	terminal := makeTerminal(ctx)

	transformer := transformer.WithMargin(paddingRows, paddingCols)
	layout := makeLayout(transformer)
	render := makeRender(transformer)

	cleaner := composite.NewCleaner(
		stack.Cleanup,
	)

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
	return wrapper_console.NewBuilder().
		Context(ctx).
		Reactive(wrapper_console.DefaultReactiveDuration).
		Color("\x1b[0;32m").
		ToTerminal()
}

func makeScreen() screen.Screen {
	landing := wrapper_screen.NewLanding()

	history := history.New(landing).ToScreen()
	pagination := pagination.New(history).
		ForceEngine(pager.EnginePage()).
		ToScreen()
	helper := help.New(pagination).ToScreen()

	return makePipeline(helper)
}

func makePipeline(scrn screen.Screen) screen.Screen {
	headerStep := wrapper_screen.NewBaseHeader()

	inlineStep := inline.InlineTransformer(
		inline.DefaultInlineSeparator,
		pipeline.NewFilter(pipeline.Tags, screen.SystemMetaTag),
		pipeline.Footer,
	)

	spacerHeader := spacer.SpacerTransformer(
		spacer.NewMeta(1, spacer.Between, pipeline.After),
		pipeline.Header,
	)

	spacerFooter := spacer.SpacerTransformer(
		spacer.NewMeta(1, spacer.Between, pipeline.Before),
		pipeline.Footer,
	)

	return pipeline.New(scrn,
		headerStep, inlineStep, spacerHeader, spacerFooter,
	).ToScreen()
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
