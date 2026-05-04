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

	nodeBase := screen_test.MockScreen{
		Name: "base",
	}.ToNode()

	stt.Stack.Push(nodeBase.Screen.Name, "lang-1", "golang")

	nodeWrapper := screen_test.MockScreen{
		Stack: nodeBase.Stack,
	}.ToNode()

	result := screen.ResultFromUIState(stt)
	result.Node = &nodeWrapper

	cleaner.Cleanup(result, stt)

	value, exists := stt.Stack.Find(nodeBase.Screen.Name, "lang-1")

	assert.True(t, exists)
	assert.Equal(t, "golang", value.Stringf())
}

func TestStack_RemovesInactiveState(t *testing.T) {
	cleaner := NewCleaner()
	stt := state.NewUIState()

	nodeBase := screen_test.MockScreen{
		Name: "base",
	}.ToNode()

	stt.Stack.Push(nodeBase.Screen.Name, "lang-1", "golang")

	nodeNext := screen_test.MockScreen{
		Name: "next",
	}.ToNode()

	nodeWrapper := screen_test.MockScreen{}.ToNode()
	nodeWrapper.Stack = nodeNext.Stack

	result := screen.ResultFromUIState(stt)
	result.Node = &nodeWrapper

	cleaner.Cleanup(result, stt)

	_, exists := stt.Stack.Find(nodeBase.Screen.Name, "lang-1")
	assert.False(t, exists)

	stt.Stack.Push(nodeNext.Screen.Name, "lang-2", "ziglang")

	value, exists := stt.Stack.Find(nodeNext.Screen.Name, "lang-2")
	assert.True(t, exists)
	assert.Equal(t, "ziglang", value.Stringf())
}

func TestStack_TransitionBetweenScreens(t *testing.T) {
	cleaner := NewCleaner()
	stt := state.NewUIState()

	nodeBase := screen_test.MockScreen{
		Name: "base",
	}.ToNode()

	nodeNext := screen_test.MockScreen{
		Name: "next",
	}.ToNode()

	stt.Stack.Push(nodeBase.Screen.Name, "lang-1", "golang")

	nodeWrapper := screen_test.MockScreen{}.ToNode()
	nodeWrapper.Stack = nodeBase.Stack

	result := screen.ResultFromUIState(stt)
	result.Node = &nodeWrapper
	cleaner.Cleanup(result, stt)

	_, exists := stt.Stack.Find(nodeBase.Screen.Name, "lang-1")
	assert.True(t, exists)

	nodeWrapper.Stack = nodeNext.Stack

	result = screen.ResultFromUIState(stt)
	result.Node = &nodeWrapper
	cleaner.Cleanup(result, stt)

	_, exists = stt.Stack.Find(nodeBase.Screen.Name, "lang-1")
	assert.False(t, exists)
}
