package wrapper_screen

import (
	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/core/screen"
	"github.com/Rafael24595/go-terminal/engine/core/screen/commons"
)

func NewBaseHeader(screen screen.Screen) screen.Screen {
	header := core.FixedLinesFromLines(
		core.ModePadding(core.Center),
		core.LineFromString("LOREM IPSUM DOLOR SIT AMET"),
		core.LineFromString("CONSECTETUR ADIPISCING"),
		core.LineFromString("-SERVER 00-"),
		core.LineJump(),
	)

	return commons.NewHeader(screen).
		AddHeader(header...).
		ToScreen()
}
