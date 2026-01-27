package wrapper_screen

import (
	"github.com/Rafael24595/go-terminal/engine/core"
	wrapper_commons "github.com/Rafael24595/go-terminal/wrapper/screen/commons"
)

func NewBaseHeader(screen core.Screen) core.Screen {
	header := core.FixedLinesFromLines(
		core.ModePadding(core.Center),
		core.LineFromString("LOREM IPSUM DOLOR SIT AMET"),
		core.LineFromString("CONSECTETUR ADIPISCING"),
		core.LineFromString("-SERVER 00-"),
		core.LineJump(),
	)

	return wrapper_commons.NewHeader(screen).
		AddHeader(header...).
		ToScreen()
}
