package wrapper_screen

import "github.com/Rafael24595/go-terminal/engine/core"

type WrapperMain struct {
}

func NewWrapperMain() *WrapperMain {
	return &WrapperMain{}
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
		headerPadding,
		"LOREM IPSUM DOLOR SIT AMET",
	))

	lines = append(lines, core.NewLine(
		headerPadding,
		"CONSECTETUR ADIPISCING",
	))

	lines = append(lines, core.NewLine(
		headerPadding,
		"-SERVER 00-",
	))

	lines = append(lines, core.NewLine(
		core.Padding{
			Padding: core.Fill,
		}, "=",
	))

	return core.ViewModel{
		Lines: lines,
	}
}
