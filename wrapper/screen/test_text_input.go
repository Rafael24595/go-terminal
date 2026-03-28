package wrapper_screen

import (
	"github.com/Rafael24595/go-terminal/engine/app/screen"
	"github.com/Rafael24595/go-terminal/engine/app/screen/primitive"
)

func NewTestTextInput() screen.Screen {
	return primitive.NewTextInput().
		SetName("textinput - amet").
		EnableBlinking().
		AddText("AD Lorem ipsum dolor sit amet.").
		ToScreen()
}
