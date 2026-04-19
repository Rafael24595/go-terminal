package mapdrw

import "github.com/Rafael24595/go-terminal/engine/terminal"

func PredFixedWinsize(size terminal.Winsize) drawInputPred {
	return func(_ terminal.Winsize) terminal.Winsize {
		return size
	}
}
