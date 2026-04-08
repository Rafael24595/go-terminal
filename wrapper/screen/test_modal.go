package wrapper_screen

import (
	"github.com/Rafael24595/go-terminal/engine/app/screen"
	"github.com/Rafael24595/go-terminal/engine/app/screen/primitive"
	"github.com/Rafael24595/go-terminal/engine/model/input"
	"github.com/Rafael24595/go-terminal/engine/render/text"
)

func NewTestModal() screen.Screen {
	return primitive.NewModalMenu().
		SetName("modal - dolor").
		AddText(
			text.NewLine("AD Lorem ipsum dolor sit amet"),
		).
		AddOptions([]input.MenuOption{
			{
				Label:  text.NewFragment("Option_1"),
				Action: NewTestModal,
			},
			{
				Label:  text.NewFragment("Option_2"),
				Action: NewTestModal,
			},
			{
				Label:  text.NewFragment("Option_3"),
				Action: NewTestModal,
			},
			{
				Label:  text.NewFragment("Option_4"),
				Action: NewTestModal,
			},
		}...).
		ToScreen()
}
