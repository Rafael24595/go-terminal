package platform

import (
	"context"
	"time"

	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
)

func ResizeReactiveEvents(ctx context.Context, drt time.Duration) <-chan winsize.Winsize {
	out := make(chan winsize.Winsize, 1)
	go listenResizeReactiveEvents(ctx, drt, out)
	return out
}

func listenResizeReactiveEvents(
	ctx context.Context,
	drt time.Duration,
	out chan winsize.Winsize,
) {
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

func ResizeProactiveEvents(ctx context.Context, drt time.Duration) <-chan winsize.Winsize {
	out := make(chan winsize.Winsize, 1)
	go listenResizeProactiveEvents(ctx, drt, out)
	return out
}

func listenResizeProactiveEvents(
	ctx context.Context,
	drt time.Duration,
	out chan winsize.Winsize,
) {
	defer close(out)

	ticker := time.NewTicker(drt)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			size, _ := Size()
			select {
			case out <- size:
			default:
			}
		}
	}
}
