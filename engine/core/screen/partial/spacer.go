package partial

import (
	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/core/drawable/line"
	"github.com/Rafael24595/go-terminal/engine/core/drawable/stack"
	"github.com/Rafael24595/go-terminal/engine/core/screen"
	"github.com/Rafael24595/go-terminal/engine/core/text"
)

type SpacerMeta struct {
	size    uint8
	between bool
}

func NewSpacerMeta(size uint8, between bool) SpacerMeta {
	return SpacerMeta{
		size:    size,
		between: between,
	}
}

type Spacer struct {
	header SpacerMeta
	footer SpacerMeta
	screen screen.Screen
}

func NewSpacer(screen screen.Screen) *Spacer {
	return &Spacer{
		screen: screen,
	}
}

func (c *Spacer) Header(header SpacerMeta) *Spacer {
	c.header = header
	return c
}

func (c *Spacer) Footer(footer SpacerMeta) *Spacer {
	c.footer = footer
	return c
}

func (c *Spacer) ToScreen() screen.Screen {
	return screen.Screen{
		Name:       c.screen.Name,
		Definition: c.screen.Definition,
		Update:     c.update,
		View:       c.view,
	}
}

func (c *Spacer) update(state *state.UIState, event screen.ScreenEvent) screen.ScreenResult {
	result := c.screen.Update(state, event)
	if result.Screen != nil {
		newScreen := NewSpacer(*result.Screen).
			Header(c.header).
			Footer(c.footer).
			ToScreen()
		result.Screen = &newScreen
	}
	return result
}

func (c *Spacer) view(state state.UIState) core.ViewModel {
	vm := c.screen.View(state)

	if c.header.size > 0 && vm.Header.Size() > 0 {
		vm = c.addHeaderStyles(vm)
	}

	if c.footer.size > 0 && vm.Footer.Size() > 0 {
		vm = c.addFooterStyles(vm)
	}

	return vm
}

func (c *Spacer) addHeaderStyles(vm core.ViewModel) core.ViewModel {
	if !c.header.between {
		spacer := line.EagerDrawableFromLines(
			c.makeSpaces(c.header.size)...,
		)

		vm.Header.Shift(spacer)
		return vm
	}

	items := vm.Header.Items()

	vm.Header = stack.NewStackDrawable()
	for _, h := range items {
		spacer := line.EagerDrawableFromLines(
			c.makeSpaces(c.header.size)...,
		)
		vm.Header.Shift(h, spacer)
	}

	return vm
}

func (c *Spacer) addFooterStyles(vm core.ViewModel) core.ViewModel {
	if !c.footer.between {
		spacer := line.EagerDrawableFromLines(
			c.makeSpaces(c.footer.size)...,
		)

		vm.Footer.Unshift(spacer)
		return vm
	}

	items := vm.Footer.Items()

	vm.Footer = stack.NewStackDrawable()
	for _, h := range items {
		spacer := line.EagerDrawableFromLines(
			c.makeSpaces(c.footer.size)...,
		)
		vm.Footer.Shift(spacer, h)
	}

	return vm
}

func (c *Spacer) makeSpaces(size uint8) []text.Line {
	spaces := make([]text.Line, size)
	for i := range spaces {
		spaces[i] = text.LineJump()
	}
	return spaces
}
