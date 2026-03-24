package state

import (
	"fmt"
	"sync"
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	
	"github.com/Rafael24595/go-terminal/engine/commons/structure/set"
)

func TestStackContext_CRUD(t *testing.T) {
	ctx := newStackContext()

	screen := "Landing"
	key := "lang"
	val := "golang"

	ctx.Push(screen, key, val)
	arg, found := ctx.Find(screen, key)

	assert.True(t, found)
	assert.Equal(t, val, arg.Stringf())

	ctx.RemoveArgument(screen, key)
	_, found = ctx.Find(screen, key)

	assert.False(t, found)

	ctx.Push(screen, "order", 1)
	ctx.RemoveScreen(screen)
	_, found = ctx.Find(screen, "order")

	assert.False(t, found)
}

func TestStackContext_RetainOnly(t *testing.T) {
	ctx := newStackContext()

	ctx.Push("Home", "a", 1)
	ctx.Push("Settings", "b", 2)
	ctx.Push("Profile", "c", 3)

	keep := set.SetFrom("Home", "Profile")
	ctx.RetainOnly(keep)

	_, found := ctx.Find("Home", "a")
	assert.True(t, found)

	_, found = ctx.Find("Settings", "b")
	assert.False(t, found)
}

func TestStackContext_Concurrency(t *testing.T) {
	ctx := newStackContext()
	const workers = 15
	var wg sync.WaitGroup
	wg.Add(workers * 2)

	for i := range workers {
		go func(id int) {
			defer wg.Done()
			ctx.Push("Screen", fmt.Sprintf("k%d", id), id)
		}(i)
	}

	for i := range workers {
		go func(id int) {
			defer wg.Done()
			ctx.Find("Screen", fmt.Sprintf("k%d", id))
		}(i)
	}

	wg.Wait()
}
