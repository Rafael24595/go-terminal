package wrapper_screen

import (
	"github.com/Rafael24595/go-terminal/engine/core/screen"
	"github.com/Rafael24595/go-terminal/engine/core/screen/wrapper"
	"github.com/Rafael24595/go-terminal/engine/core/style"
	"github.com/Rafael24595/go-terminal/engine/core/text"
)

func NewBaseHeader(screen screen.Screen) screen.Screen {
	header := text.FixedLinesFromLines(
		style.SpecFromKind(style.SpcKindPaddingCenter),
		text.LineFromString("Lorem ipsum dolor sit amet", style.AtmUpper),
		text.LineFromString("consectetur adipiscing", style.AtmUpper),
		text.LineFromString("-Server 00-", style.AtmUpper),
	)

	return wrapper.NewHeader(screen).
		AddHeader(header...).
		ToScreen()
}
