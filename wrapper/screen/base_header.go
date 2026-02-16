package wrapper_screen

import (
	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/core/screen"
	"github.com/Rafael24595/go-terminal/engine/core/screen/commons"
)

func NewBaseHeader(screen screen.Screen) screen.Screen {
	header := core.FixedLinesFromLines(
		core.ModePadding(core.Center),
		core.LineFromString("Lorem ipsum dolor sit amet", core.Upper),
		core.LineFromString("consectetur adipiscing", core.Upper),
		core.LineFromString("-Server 00-", core.Upper),
		core.LineJump(),
	)

	return commons.NewHeader(screen).
		AddHeader(header...).
		ToScreen()
}
