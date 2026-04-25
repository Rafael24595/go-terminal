package pipeline

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
)

type Transformer func(viewmodel.ViewModel) viewmodel.ViewModel

type Pipeline struct {
	screen screen.Screen
	steps  []Transformer
}

func NewPipeline(screen screen.Screen, steps ...Transformer) *Pipeline {
	return &Pipeline{
		screen: screen,
		steps:  make([]Transformer, 0),
	}
}

func (c *Pipeline) PushSteps(steps ...Transformer) *Pipeline {
	c.steps = append(c.steps, steps...)
	return c
}

func (c *Pipeline) ToScreen() screen.Screen {
	return screen.Screen{
		Name:       c.screen.Name,
		Definition: c.screen.Definition,
		Update:     c.update,
		View:       c.view,
		Stack:      c.screen.Stack,
	}
}

func (c *Pipeline) update(state *state.UIState, event screen.ScreenEvent) screen.ScreenResult {
	result := c.screen.Update(state, event)
	if result.Screen != nil {
		newScreen := NewPipeline(*result.Screen).
			PushSteps(c.steps...).
			ToScreen()
		result.Screen = &newScreen
	}
	return result
}

func (c *Pipeline) view(state state.UIState) viewmodel.ViewModel {
	vm := c.screen.View(state)
	for _, s := range c.steps {
		vm = s(vm)
	}
	return vm
}
