package composite

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"

	cleaner_test "github.com/Rafael24595/go-reacterm-core/test/engine/app/cleaner"
)

func TestComposite_ToStateCleaner(t *testing.T) {
	cleaner := NewCleaner()

	cleaner_test.Helper_ToStateCleaner(t, cleaner)
}

func TestCompositeCleanup(t *testing.T) {

	c := NewCleaner(
		func(r screen.ScreenResult, s *state.UIState) *state.UIState {
			s.Pager.ActualPage = 0
			return s
		},
		func(r screen.ScreenResult, s *state.UIState) *state.UIState {
			s.Pager.ForceShow = true
			return s
		},
		func(r screen.ScreenResult, s *state.UIState) *state.UIState {
			s.Helper.ShowHelp = false
			return s
		},
	)

	stt := state.NewUIState()
	stt.Pager.ActualPage = 10
	stt.Pager.ForceShow = false
	stt.Helper.ShowHelp = true

	res := screen.ScreenResult{}

	stt = c.Cleanup(res, stt)

	assert.Equal(t, 0, stt.Pager.ActualPage)
	assert.True(t, stt.Pager.ForceShow)
	assert.False(t, stt.Helper.ShowHelp)
}
