package action

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
)

func TestMergeFocus(t *testing.T) {
	f := MergeFocus(FocusHeader, FocusBody)

	assert.Equal(t, FocusHeader|FocusBody, f)
}

func TestFocus_HasAny(t *testing.T) {
	f := MergeFocus(FocusHeader, FocusBody)

	assert.True(t, f.HasAny(FocusHeader))
	assert.True(t, f.HasAny(FocusBody))
	assert.False(t, f.HasAny(FocusFooter))
}

func TestFocus_HasNone(t *testing.T) {
	f := MergeFocus(FocusHeader)

	assert.True(t, f.HasNone(FocusBody, FocusFooter))
	assert.False(t, f.HasNone(FocusHeader))
}
