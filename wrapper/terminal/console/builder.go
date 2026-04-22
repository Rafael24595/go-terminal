package wrapper_console

import (
	"context"
	"time"

	"github.com/Rafael24595/go-terminal/engine/model/winsize"
	"github.com/Rafael24595/go-terminal/engine/terminal"
	"github.com/Rafael24595/go-terminal/wrapper/platform"
	wrapper_reader "github.com/Rafael24595/go-terminal/wrapper/terminal/reader"
)

type ConsoleBuilder struct {
	context  context.Context
	strategy resizeStrategy
	reader   *wrapper_reader.KeyReader
	color    string
}

func NewBuilder() *ConsoleBuilder {
	return &ConsoleBuilder{
		context:  context.Background(),
		strategy: defaultStrategy(),
		reader:   wrapper_reader.New(),
		color:    "",
	}
}

func (b *ConsoleBuilder) Context(context context.Context) *ConsoleBuilder {
	b.context = context
	return b
}

func (b *ConsoleBuilder) Reactive(duration time.Duration) *ConsoleBuilder {
	b.strategy = func(ctx context.Context) <-chan winsize.Winsize {
		return platform.ResizeSystemEvents(ctx, duration)
	}
	return b
}

func (b *ConsoleBuilder) Proactive(duration time.Duration) *ConsoleBuilder {
	b.strategy = func(ctx context.Context) <-chan winsize.Winsize {
		return platform.ResizeProactiveEvents(ctx, duration)
	}
	return b
}

func (b *ConsoleBuilder) Color(color string) *ConsoleBuilder {
	b.color = color
	return b
}

func (b *ConsoleBuilder) Build() *Console {
	console := newConsole()
	console.context = b.context
	console.strategy = b.strategy
	console.color = b.color

	return console
}

func (b *ConsoleBuilder) ToTerminal() terminal.Terminal {
	return b.Build().ToTerminal()
}
