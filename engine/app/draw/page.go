package draw

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
)

type PageRenderer func(*state.UIState, winsize.Winsize, drawable.Unit) *DrawState
