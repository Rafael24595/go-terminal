package action_test

import (
	"testing"

	"github.com/Rafael24595/go-terminal/engine/core/drawable"
	"github.com/Rafael24595/go-terminal/engine/core/drawable/action"
	"github.com/Rafael24595/go-terminal/test/support/assert"
)

func TestMergeFocus(t *testing.T) {
	f := action.MergeFocus(action.FocusHeader, action.FocusBody)

	assert.Equal(t, action.FocusHeader|action.FocusBody, f)
}

func TestFocus_HasAny(t *testing.T) {
	f := action.MergeFocus(action.FocusHeader, action.FocusBody)

	assert.True(t, f.HasAny(action.FocusHeader))
	assert.True(t, f.HasAny(action.FocusBody))
	assert.False(t, f.HasAny(action.FocusFooter))
}

func TestFocus_HasNone(t *testing.T) {
	f := action.MergeFocus(action.FocusHeader)

	assert.True(t, f.HasNone(action.FocusBody, action.FocusFooter))
	assert.False(t, f.HasNone(action.FocusHeader))
}

func TestApplyAction_Map(t *testing.T) {
	count := 0

	act := action.NewAction(action.ActionMapEach, action.FocusNone,
		func(items ...drawable.Drawable) []drawable.Drawable {
			count++
			return items
		},
	)

	drw := drawable.Drawable{}

	action.ApplyAction(act, drw, drw)

	assert.Equal(t, 2, count)
}

func TestApplyAction_Group(t *testing.T) {
	count := 0
	size := 0

	act := action.NewAction(action.ActionMapGroup, action.FocusNone,
		func(items ...drawable.Drawable) []drawable.Drawable {
			count++
			size = len(items)
			return items
		},
	)

	drw := drawable.Drawable{}

	action.ApplyAction(act, drw, drw, drw)

	assert.Equal(t, 1, count)
	assert.Equal(t, 3, size)
}
