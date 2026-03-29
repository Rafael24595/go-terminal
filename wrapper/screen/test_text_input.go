package wrapper_screen

import (
	"github.com/Rafael24595/go-terminal/engine/app/screen"
	"github.com/Rafael24595/go-terminal/engine/app/screen/primitive"
	"github.com/Rafael24595/go-terminal/engine/render/text"
)

func NewTestTextInput() screen.Screen {
	label := text.FragmentsFromString("Ipsum")
	return primitive.NewTextInput().
		SetName("textinput - amet").
		EnableBlinking().
		SetLabel(label).
		AddText("AD Lorem ipsum dolor sit amet.").
		ToScreen()
}
