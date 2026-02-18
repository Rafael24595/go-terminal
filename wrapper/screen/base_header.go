package wrapper_screen

import (
	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/core/screen"
	"github.com/Rafael24595/go-terminal/engine/core/screen/commons"
	"github.com/Rafael24595/go-terminal/engine/core/style"
)

func NewBaseHeader(screen screen.Screen) screen.Screen {
	header := core.FixedLinesFromLines(
		style.SpecFromKind(style.SpcKindCenter),
		core.LineFromString("Lorem ipsum dolor sit amet", style.AtmUpper),
		core.LineFromString("consectetur adipiscing", style.AtmUpper),
		core.LineFromString("-Server 00-", style.AtmUpper),
		core.LineJump(),
	)

	return commons.NewHeader(screen).
		AddHeader(header...).
		ToScreen()
}
