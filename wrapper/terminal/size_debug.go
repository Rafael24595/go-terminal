//go:build mock_cmd
// +build mock_cmd

package wrapper_console

import (
	core_terminal "github.com/Rafael24595/go-terminal/engine/terminal"
)

func Size() core_terminal.Winsize {
	return core_terminal.Winsize{
		Rows: 80,
		Cols: 150,
	}
}
