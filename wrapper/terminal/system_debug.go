//go:build mock_cmd
// +build mock_cmd

package wrapper_terminal

import (
	"context"
	"time"

	"github.com/Rafael24595/go-terminal/engine/model/winsize"
)

func onStart() (uintptr, error) {
	return 0, nil
}

func onClose(rawmode uintptr) {}

func Size() core_winsize.Winsize {
	return core_winsize.New(80, 150)
}

func ResizeEvents(ctx context.Context) <-chan winsize.Winsize {
	return timeResizeSignal(ctx, 10*time.Millisecond)
}
