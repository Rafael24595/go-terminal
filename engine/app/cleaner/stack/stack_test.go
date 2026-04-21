package stack

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-terminal/engine/app/screen"
	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/commons/structure/set"

	cleaner_test "github.com/Rafael24595/go-terminal/test/engine/app/cleaner"
	screen_test "github.com/Rafael24595/go-terminal/test/engine/app/screen"
)

func TestStack_ToStateCleaner(t *testing.T) {
	cleaner := NewCleaner()

	cleaner_test.Helper_ToStateCleaner(t, cleaner)
}

func TestStack_PreservesActiveState(t *testing.T) {
	cleaner := NewCleaner()
	stt := state.NewUIState()

	screenBase := screen_test.MockScreen{
		Name: "base",
	}.ToScreen().StackFromName()

	stt.Stack.Push(screenBase.Name(), "lang-1", "golang")

	screenWrapper := screen_test.MockScreen{}.ToScreen()
	screenWrapper.Stack = func() set.Set[string] {
		return screenBase.Stack()
	}

	result := screen.ScreenResultFromUIState(stt)
	result.Screen = &screenWrapper

	cleaner.Cleanup(result, stt)

	value, exists := stt.Stack.Find(screenBase.Name(), "lang-1")

	assert.True(t, exists)
	assert.Equal(t, "golang", value.Stringf())
}

func TestStack_RemovesInactiveState(t *testing.T) {
	cleaner := NewCleaner()
	stt := state.NewUIState()

	screenBase := screen_test.MockScreen{
		Name: "base",
	}.ToScreen().StackFromName()

	stt.Stack.Push(screenBase.Name(), "lang-1", "golang")

	screenNext := screen_test.MockScreen{
		Name: "next",
	}.ToScreen().StackFromName()

	screenWrapper := screen_test.MockScreen{}.ToScreen()
	screenWrapper.Stack = func() set.Set[string] {
		return screenNext.Stack()
	}

	result := screen.ScreenResultFromUIState(stt)
	result.Screen = &screenWrapper

	cleaner.Cleanup(result, stt)

	_, exists := stt.Stack.Find(screenBase.Name(), "lang-1")
	assert.False(t, exists)

	stt.Stack.Push(screenNext.Name(), "lang-2", "ziglang")

	value, exists := stt.Stack.Find(screenNext.Name(), "lang-2")
	assert.True(t, exists)
	assert.Equal(t, "ziglang", value.Stringf())
}

func TestStack_TransitionBetweenScreens(t *testing.T) {
	cleaner := NewCleaner()
	stt := state.NewUIState()

	screenBase := screen_test.MockScreen{
		Name: "base",
	}.ToScreen().StackFromName()

	screenNext := screen_test.MockScreen{
		Name: "next",
	}.ToScreen().StackFromName()

	stt.Stack.Push(screenBase.Name(), "lang-1", "golang")

	screenWrapper := screen_test.MockScreen{}.ToScreen()
	screenWrapper.Stack = func() set.Set[string] {
		return screenBase.Stack()
	}

	result := screen.ScreenResultFromUIState(stt)
	result.Screen = &screenWrapper
	cleaner.Cleanup(result, stt)

	_, exists := stt.Stack.Find(screenBase.Name(), "lang-1")
	assert.True(t, exists)

	screenWrapper.Stack = func() set.Set[string] {
		return screenNext.Stack()
	}

	result = screen.ScreenResultFromUIState(stt)
	result.Screen = &screenWrapper
	cleaner.Cleanup(result, stt)

	_, exists = stt.Stack.Find(screenBase.Name(), "lang-1")
	assert.False(t, exists)
}
