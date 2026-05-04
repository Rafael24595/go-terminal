package stack

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"

	cleaner_test "github.com/Rafael24595/go-reacterm-core/test/engine/app/cleaner"
	screen_test "github.com/Rafael24595/go-reacterm-core/test/engine/app/screen"
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
	}.ToScreen()

	stt.Stack.Push(screenBase.Name, "lang-1", "golang")

	screenWrapper := screen_test.MockScreen{
		Stack: screenBase.Stack,
	}.ToScreen()

	result := screen.ResultFromUIState(stt)
	result.Screen = &screenWrapper

	cleaner.Cleanup(result, stt)

	value, exists := stt.Stack.Find(screenBase.Name, "lang-1")

	assert.True(t, exists)
	assert.Equal(t, "golang", value.Stringf())
}

func TestStack_RemovesInactiveState(t *testing.T) {
	cleaner := NewCleaner()
	stt := state.NewUIState()

	screenBase := screen_test.MockScreen{
		Name: "base",
	}.ToScreen()

	stt.Stack.Push(screenBase.Name, "lang-1", "golang")

	screenNext := screen_test.MockScreen{
		Name: "next",
	}.ToScreen()

	screenWrapper := screen_test.MockScreen{}.ToScreen()
	screenWrapper.Stack = screenNext.Stack

	result := screen.ResultFromUIState(stt)
	result.Screen = &screenWrapper

	cleaner.Cleanup(result, stt)

	_, exists := stt.Stack.Find(screenBase.Name, "lang-1")
	assert.False(t, exists)

	stt.Stack.Push(screenNext.Name, "lang-2", "ziglang")

	value, exists := stt.Stack.Find(screenNext.Name, "lang-2")
	assert.True(t, exists)
	assert.Equal(t, "ziglang", value.Stringf())
}

func TestStack_TransitionBetweenScreens(t *testing.T) {
	cleaner := NewCleaner()
	stt := state.NewUIState()

	screenBase := screen_test.MockScreen{
		Name: "base",
	}.ToScreen()

	screenNext := screen_test.MockScreen{
		Name: "next",
	}.ToScreen()

	stt.Stack.Push(screenBase.Name, "lang-1", "golang")

	screenWrapper := screen_test.MockScreen{}.ToScreen()
	screenWrapper.Stack = screenBase.Stack

	result := screen.ResultFromUIState(stt)
	result.Screen = &screenWrapper
	cleaner.Cleanup(result, stt)

	_, exists := stt.Stack.Find(screenBase.Name, "lang-1")
	assert.True(t, exists)

	screenWrapper.Stack = screenNext.Stack

	result = screen.ResultFromUIState(stt)
	result.Screen = &screenWrapper
	cleaner.Cleanup(result, stt)

	_, exists = stt.Stack.Find(screenBase.Name, "lang-1")
	assert.False(t, exists)
}
