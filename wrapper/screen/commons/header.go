package wrapper_commons

import "github.com/Rafael24595/go-terminal/engine/core"

type Header struct {
	screen core.Screen
	header []core.Line
}

func NewHeader(screen core.Screen) *Header {
	return &Header{
		screen: screen,
		header: make([]core.Line, 0),
	}
}

func (c *Header) AddHeader(header ...core.Line) *Header {
	c.header = append(c.header, header...)
	return c
}

func (c *Header) ToScreen() core.Screen {
	return core.Screen{
		Update: c.Update,
		View:   c.View,
	}
}

func (c *Header) Update(e core.ScreenEvent) {
	c.screen.Update(e)
}

func (c *Header) View() core.ViewModel {
	vm := c.screen.View()
	return core.ViewModel{
		Headers: c.header,
		Lines:   vm.Lines,
		Input:   vm.Input,
	}
}
