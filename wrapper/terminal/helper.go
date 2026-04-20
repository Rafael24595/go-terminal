package wrapper_terminal

import (
	"context"
	"time"

	"github.com/Rafael24595/go-terminal/engine/model/winsize"
)

func timeResizeEvents(ctx context.Context, drt time.Duration) <-chan winsize.Winsize {
	out := make(chan winsize.Winsize, 1)
	go listenTimeResizeEvents(ctx, drt, out)
	return out
}

func listenTimeResizeEvents(ctx context.Context, drt time.Duration, out chan winsize.Winsize) {
	defer close(out)

	ticker := time.NewTicker(drt)
	defer ticker.Stop()

	lastSize, _ := Size()
	
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			size, err := Size()
			if err != nil || size.Eq(lastSize) {
				continue
			}

			lastSize = size

			select {
			case out <- size:
			default:
			}
		}
	}
}
