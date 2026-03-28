package partial

import (
	"github.com/Rafael24595/go-terminal/engine/app/screen"
	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/app/viewmodel"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable/line"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable/stack"
	"github.com/Rafael24595/go-terminal/engine/render/spacer"
	"github.com/Rafael24595/go-terminal/engine/render/text"
)

type Spacer struct {
	header spacer.SpacerMeta
	footer spacer.SpacerMeta
	screen screen.Screen
}

func NewSpacer(screen screen.Screen) *Spacer {
	return &Spacer{
		screen: screen,
	}
}

func (c *Spacer) Header(header spacer.SpacerMeta) *Spacer {
	c.header = header
	return c
}

func (c *Spacer) Footer(footer spacer.SpacerMeta) *Spacer {
	c.footer = footer
	return c
}

func (c *Spacer) ToScreen() screen.Screen {
	return screen.Screen{
		Name:       c.screen.Name,
		Definition: c.screen.Definition,
		Update:     c.update,
		View:       c.view,
		Stack:      c.screen.Stack,
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

func (c *Spacer) view(state state.UIState) viewmodel.ViewModel {
	vm := c.screen.View(state)

	if c.header.Size > 0 && vm.Header.Size() > 0 {
		vm = c.addHeaderStyles(vm)
	}

	if c.footer.Size > 0 && vm.Footer.Size() > 0 {
		vm = c.addFooterStyles(vm)
	}

	return vm
}

func (c *Spacer) addHeaderStyles(vm viewmodel.ViewModel) viewmodel.ViewModel {
	spcr := line.EagerDrawableFromLines(
		c.makeSpaces(c.header.Size)...,
	)

	if c.header.Mode == spacer.SpacerAppend {
		vm.Header.Push(spcr)
		return vm
	}

	items := vm.Header.Items()

	vm.Header = stack.NewStackDrawable()
	for _, h := range items {
		vm.Header.Push(h, spcr)
	}

	return vm
}

func (c *Spacer) addFooterStyles(vm viewmodel.ViewModel) viewmodel.ViewModel {
	spcr := line.EagerDrawableFromLines(
		c.makeSpaces(c.footer.Size)...,
	)

	if c.footer.Mode == spacer.SpacerAppend {
		vm.Footer.Unshift(spcr)
		return vm
	}

	items := vm.Footer.Items()

	vm.Footer = stack.NewStackDrawable()
	for _, h := range items {
		vm.Footer.Push(spcr, h)
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
