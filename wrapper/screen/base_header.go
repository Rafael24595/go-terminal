package wrapper_screen

import (
	"github.com/Rafael24595/go-terminal/engine/app/screen"
	"github.com/Rafael24595/go-terminal/engine/app/screen/wrapper"
	"github.com/Rafael24595/go-terminal/engine/render/style"
	"github.com/Rafael24595/go-terminal/engine/render/text"
)

func NewBaseHeader(screen screen.Screen) screen.Screen {
	header := text.ApplyLineSpec(
		style.SpecFromKind(style.SpcKindPaddingCenter),
		*text.LineFromFragments(
			*text.NewFragment("Lorem ipsum dolor sit amet").AddAtom(style.AtmUpper),
		),
		*text.LineFromFragments(
			*text.NewFragment("consectetur adipiscing").AddAtom(style.AtmUpper),
		),
		*text.LineFromFragments(
			*text.NewFragment("-Server 00-").AddAtom(style.AtmUpper),
		),
	)

	return wrapper.NewHeader(screen).
		AddHeader(header...).
		ToScreen()
}
