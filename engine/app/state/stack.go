package state

import (
	"sync"

	"github.com/Rafael24595/go-terminal/engine/commons"
	"github.com/Rafael24595/go-terminal/engine/commons/structure/set"
	"github.com/Rafael24595/go-terminal/engine/model/param"
	"github.com/Rafael24595/go-terminal/engine/platform/clock"
)

type StackContext struct {
	mu      sync.RWMutex
	clock   clock.Clock
	context map[string]*ScreenContext
}

func newStackContext() *StackContext {
	return &StackContext{
		clock:   clock.UnixMilliClock,
		context: make(map[string]*ScreenContext),
	}
}

func (c *StackContext) Find(screen, key string) (*commons.Argument, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	ctx, ok := c.context[screen]
	if !ok {
		return nil, false
	}

	return ctx.Find(key)
}

func (c *StackContext) Push(screen, key string, arg any) *StackContext {
	c.mu.Lock()
	defer c.mu.Unlock()

	ctx, ok := c.context[screen]
	if !ok {
		ctx = newScreenContext(c.clock)
	}

	c.context[screen] = ctx.Push(key,
		newContextArgument(c.clock, arg),
	)

	return c
}

func (c *StackContext) RemoveScreen(screen string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, ok := c.context[screen]
	if !ok {
		return false
	}

	delete(c.context, screen)

	return true
}

func (c *StackContext) RemoveArgument(screen, key string) (*commons.Argument, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	ctx, ok := c.context[screen]
	if !ok {
		return nil, false
	}

	return ctx.Remove(key)
}

func (c *StackContext) RetainOnly(screens set.Set[string]) *StackContext {
	c.mu.Lock()
	items := make([]string, 0)
	for screen := range c.context {
		if !screens.Has(screen) {
			items = append(items, screen)
		}
	}
	c.mu.Unlock()

	for _, name := range items {
		c.RemoveScreen(name)
	}

	return c
}

func PushParam[T any](
	c *StackContext,
	screen string,
	param param.Typed[T],
	arg T,
) *StackContext {
	return c.Push(screen, param.Code(), arg)
}

type ScreenContext struct {
	mu        sync.RWMutex
	timestamp int64
	context   map[string]ContextArgument
}

func newScreenContext(clock clock.Clock) *ScreenContext {
	return &ScreenContext{
		timestamp: clock(),
		context:   make(map[string]ContextArgument),
	}
}

func (c *ScreenContext) Find(key string) (*commons.Argument, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	arg, ok := c.context[key]
	if !ok {
		return nil, false
	}

	return &arg.argument, true
}

func (c *ScreenContext) Push(key string, arg ContextArgument) *ScreenContext {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.context[key] = arg

	return c
}

func (c *ScreenContext) Remove(key string) (*commons.Argument, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	arg, ok := c.context[key]
	if !ok {
		return nil, false
	}

	delete(c.context, key)

	return &arg.argument, true
}

type ContextArgument struct {
	clock     clock.Clock
	timestamp int64
	argument  commons.Argument
}

func newContextArgument(clk clock.Clock, arg any) ContextArgument {
	return ContextArgument{
		clock:     clk,
		timestamp: clk(),
		argument:  commons.ArgumentFrom(arg),
	}
}
