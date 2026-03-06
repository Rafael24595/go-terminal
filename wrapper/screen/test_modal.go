package wrapper_screen

import (
	"github.com/Rafael24595/go-terminal/engine/core/input"
	"github.com/Rafael24595/go-terminal/engine/core/screen"
	"github.com/Rafael24595/go-terminal/engine/core/screen/primitive"
	"github.com/Rafael24595/go-terminal/engine/core/text"
)

func NewTestModal() screen.Screen {
	return primitive.NewModalMenu().
		SetName("modal - dolor").
		AddText(
			text.LineFromString("AD Lorem ipsum dolor sit amet"),
		).
		AddOptions([]input.MenuOption{
			{
				Label: text.NewFragment("Option_1"),
				Action: NewTestModal,
			},
			{
				Label: text.NewFragment("Option_2"),
				Action: NewTestModal,
			},
			{
				Label: text.NewFragment("Option_3"),
				Action: NewTestModal,
			},
			{
				Label: text.NewFragment("Option_4"),
				Action: NewTestModal,
			},
		}...).
		ToScreen()
}
