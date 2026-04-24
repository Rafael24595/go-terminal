package wrapper_screen

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/primitive"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
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
