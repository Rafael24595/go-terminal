package wrapper_console

import (
	"context"
	"time"

	"github.com/Rafael24595/go-terminal/engine/model/winsize"
	"github.com/Rafael24595/go-terminal/wrapper/platform"
)

const DefaultReactiveDuration = 150 * time.Millisecond

type resizeStrategy func(ctx context.Context) <-chan winsize.Winsize

var defaultStrategy = func() resizeStrategy {
	return func(ctx context.Context) <-chan winsize.Winsize {
		return platform.ResizeSystemEvents(ctx, DefaultReactiveDuration)
	}
}
