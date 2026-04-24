//go:build mock_cmd
// +build mock_cmd

package platform

import (
	"context"
	"time"

	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
)

func OnStart() (uintptr, error) {
	return 0, nil
}

func OnClose(rawmode uintptr) {}

func Size() core_winsize.Winsize {
	return core_winsize.New(80, 150)
}

func ResizeSystemEvents(ctx context.Context, _ time.Duration) <-chan winsize.Winsize {
	return ResizeReactiveEvents(ctx, 10*time.Millisecond)
}
