package wrapper_screen

import (
	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/core/screen"
	"github.com/Rafael24595/go-terminal/engine/core/screen/commons"
)

func NewTestTextArea() screen.Screen {
	return commons.NewTextArea().
		SetName("article - amet").
		AddTitle(
			core.LineFromString("Suspendisse sem arcu"),
			core.NewLine("=", core.ModePadding(core.FillUp)),
			core.LineJump(),
		).
		ToScreen()
}
