package wrapper_screen

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
	
	text_screen "github.com/Rafael24595/go-reacterm-core/engine/app/screen/primitive/text"
)

func NewTestTextInput() screen.Screen {
	label := text.FragmentsFromString("Ipsum")
	return text_screen.NewInput().
		SetName("textinput - amet").
		EnableBlinking().
		SetLabel(label).
		AddText("AD Lorem ipsum dolor sit amet.").
		ToScreen()
}
