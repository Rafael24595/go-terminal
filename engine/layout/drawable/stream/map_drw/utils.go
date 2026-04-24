package mapdrw

import "github.com/Rafael24595/go-reacterm-core/engine/model/winsize"

func PredFixedWinsize(size winsize.Winsize) drawInputPred {
	return func(_ winsize.Winsize) winsize.Winsize {
		return size
	}
}
