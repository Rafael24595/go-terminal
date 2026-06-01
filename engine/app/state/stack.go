package state

import (
	"sync"

	"github.com/Rafael24595/go-reacterm-core/engine/commons"
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/set"
	"github.com/Rafael24595/go-reacterm-core/engine/model/param"
	"github.com/Rafael24595/go-reacterm-core/engine/platform/clock"
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

func (n *StackContext) Find(screen, key string) (*commons.Argument, bool) {
	n.mu.RLock()
	defer n.mu.RUnlock()

	ctx, ok := n.context[screen]
	if !ok {
		return nil, false
	}

	return ctx.Find(key)
}

func (n *StackContext) Push(screen, key string, arg any) *StackContext {
	n.mu.Lock()
	defer n.mu.Unlock()

	ctx, ok := n.context[screen]
	if !ok {
		ctx = newScreenContext(n.clock)
	}

	n.context[screen] = ctx.Push(key,
		newContextArgument(n.clock, arg),
	)

	return n
}

func (n *StackContext) RemoveScreen(screen string) bool {
	n.mu.Lock()
	defer n.mu.Unlock()

	_, ok := n.context[screen]
	if !ok {
		return false
	}

	delete(n.context, screen)

	return true
}

func (n *StackContext) RemoveArgument(screen, key string) (*commons.Argument, bool) {
	n.mu.Lock()
	defer n.mu.Unlock()

	ctx, ok := n.context[screen]
	if !ok {
		return nil, false
	}

	return ctx.Remove(key)
}

func (n *StackContext) RetainOnly(screens set.Set[string]) *StackContext {
	n.mu.Lock()
	items := make([]string, 0)
	for screen := range n.context {
		if !screens.Has(screen) {
			items = append(items, screen)
		}
	}
	n.mu.Unlock()

	for _, name := range items {
		n.RemoveScreen(name)
	}

	return n
}

func PushParam[T any](
	c *StackContext,
	screen string,
	param param.Typed[T],
	arg T,
) *StackContext {
	return c.Push(screen, param.Code(), arg)
}

func FindParam[T any](
	c *StackContext,
	screen string,
	param param.Typed[T],
) (T, bool) {
	arg, ok := c.Find(screen, param.Code())
	if !ok {
		var zero T
		return zero, false
	}

	return commons.Map[T](*arg)
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

func (n *ScreenContext) Find(key string) (*commons.Argument, bool) {
	n.mu.RLock()
	defer n.mu.RUnlock()

	arg, ok := n.context[key]
	if !ok {
		return nil, false
	}

	return &arg.argument, true
}

func (n *ScreenContext) Push(key string, arg ContextArgument) *ScreenContext {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.context[key] = arg

	return n
}

func (n *ScreenContext) Remove(key string) (*commons.Argument, bool) {
	n.mu.Lock()
	defer n.mu.Unlock()

	arg, ok := n.context[key]
	if !ok {
		return nil, false
	}

	delete(n.context, key)

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
