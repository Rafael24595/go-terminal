package pager

import (
	"testing"

	"github.com/Rafael24595/go-terminal/engine/app/state"

	assert "github.com/Rafael24595/go-assert/assert/test"
)

func TestPagerStrategy_Integration(t *testing.T) {
	strategy := NewStrategy()

	ctx := PredicateContext{
		Page:     1,
		HasFocus: false,
	}

	state := state.UIState{
		Pager: state.PagerContext{
			TargetPage: 1,
		},
	}

	shouldStop := strategy.Predicate.Func(state, ctx)

	assert.True(t, shouldStop)
}
