//go:build mock_cmd
// +build mock_cmd

package wrapper_console

import (
	core_terminal "github.com/Rafael24595/go-terminal/engine/terminal"
)

func onStart() (uintptr, error) {
	return 0, nil
}

func onClose(rawmode uintptr) {}

func Size() core_terminal.Winsize {
	return core_terminal.Winsize{
		Rows: 80,
		Cols: 150,
	}
}

