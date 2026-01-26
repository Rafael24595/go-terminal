package wrapper_screen

import "github.com/Rafael24595/go-terminal/engine/core"

type WrapperMain struct {
	Screen core.Screen
}

func NewWrapperMain(screen core.Screen) *WrapperMain {
	return &WrapperMain{
		Screen: screen,
	}
}

func (c *WrapperMain) ToScreen() core.Screen {
	return core.Screen{
		View: c.View,
	}
}

func (c *WrapperMain) View() core.ViewModel {
	lines := make([]core.Line, 0)

	headerPadding := core.Padding{
		Padding: core.Center,
	}

	lines = append(lines, core.NewLine(
		"LOREM IPSUM DOLOR SIT AMET",
		headerPadding,
	))

	lines = append(lines, core.NewLine(
		"CONSECTETUR ADIPISCING",
		headerPadding,
	))

	lines = append(lines, core.NewLine(
		"-SERVER 00-",
		headerPadding,
	))

	lines = append(lines, core.NewLine(
		"",
		core.Padding{
			Padding: core.Fill,
		},
	))
	
	vm := c.Screen.View()

	return core.ViewModel{
		Headers: lines,
		Lines:   vm.Lines,
		Input:   vm.Input,
	}
}
