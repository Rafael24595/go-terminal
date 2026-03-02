package wrapper_screen

import (
	"github.com/Rafael24595/go-terminal/engine/core/screen"
	"github.com/Rafael24595/go-terminal/engine/core/screen/commons"
	"github.com/Rafael24595/go-terminal/engine/core/text"
)

func NewTestModal() screen.Screen {
	return commons.NewModalMenu().
		SetName("modal - dolor").
		AddOptions([]commons.ModalOption{
			{
				Fragment: text.NewFragment("Option_1"),
				Action: NewTestModal,
			},
			{
				Fragment: text.NewFragment("Option_2"),
				Action: NewTestModal,
			},
		}...).
		ToScreen()
}
