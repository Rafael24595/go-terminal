package wrapper_screen

import (
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
		AddOptions([]primitive.ModalOption{
			{
				Fragment: text.NewFragment("Option_1"),
				Action: NewTestModal,
			},
			{
				Fragment: text.NewFragment("Option_2"),
				Action: NewTestModal,
			},
			{
				Fragment: text.NewFragment("Option_3"),
				Action: NewTestModal,
			},
			{
				Fragment: text.NewFragment("Option_4"),
				Action: NewTestModal,
			},
		}...).
		ToScreen()
}
