package wrapper_screen

import "github.com/Rafael24595/go-terminal/engine/core"

type Index struct {
}

func NewIndex() *Index {
	return &Index{}
}

func (c *Index) ToScreen() core.Screen {
	return core.Screen{
		View: c.View,
	}
}

func (c *Index) View() core.ViewModel {
	lines := make([]core.Line, 0)

	lines = append(lines, core.NewLine(
		"Sed facilisis, leo sit amet molestie congue, justo risus bibendum tortor",
		core.Padding{
			Padding: core.Right,
		},
	))

	lines = append(lines, core.NewLine(
		"-",
		core.Padding{
			Padding: core.FillUp,
		},
	))

	lines = append(lines, core.NewLine(
		"- Option 0",
		core.Padding{
			Padding: core.Custom,
			Left: 2,
		},
	))

	lines = append(lines, core.NewLine(
		"- Option 1",
		core.Padding{
			Padding: core.Custom,
			Left: 2,
		},
	))

	lines = append(lines, core.NewLine(
		"- Option 2",
		core.Padding{
			Padding: core.Custom,
			Left: 2,
		},
	))

	return core.ViewModel{
		Lines: lines,
	}
}
