package pager

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	"github.com/Rafael24595/go-terminal/engine/app/state"
)

func TestPredicatePage(t *testing.T) {
	p := PredicatePage()

	state := state.UIState{
		Pager: state.PagerState{
			Page: 2,
		},
	}

	tests := []struct {
		name string
		ctx  PredicateContext
		want bool
	}{
		{
			name: "same page",
			ctx:  PredicateContext{Page: 2},
			want: true,
		},
		{
			name: "different page",
			ctx:  PredicateContext{Page: 1},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := p.Func(state, tt.ctx)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPredicateFocus(t *testing.T) {
	p := PredicateFocus()

	tests := []struct {
		name string
		ctx  PredicateContext
		want bool
	}{
		{
			name: "has focus",
			ctx:  PredicateContext{HasFocus: true},
			want: true,
		},
		{
			name: "no focus",
			ctx:  PredicateContext{HasFocus: false},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := p.Func(state.UIState{}, tt.ctx)
			assert.Equal(t, tt.want, got)
		})
	}
}
