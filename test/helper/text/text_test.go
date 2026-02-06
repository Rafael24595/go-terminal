package text_test

import (
	"testing"

	"github.com/Rafael24595/go-terminal/engine/core/key"
	"github.com/Rafael24595/go-terminal/engine/helper/text"
	"github.com/Rafael24595/go-terminal/test/support/assert"
)

func TestAddSpaceAfterRunes_AddsSpace(t *testing.T) {
	start, end, text, ok := text.AppendSpaceAfterRune(
		key.Key{Rune: ','},
		5,
		5,
		nil,
	)

	assert.True(t, ok)
	assert.Equal(t, uint(5), start)
	assert.Equal(t, uint(5), end)
	assert.Equal(t, ", ", string(text))
}

func TestAddSpaceAfterRunes_IgnoresOtherRunes(t *testing.T) {
	_, _, text, ok := text.AppendSpaceAfterRune(
		key.Key{Rune: 'a'},
		1,
		1,
		nil,
	)

	assert.False(t, ok)
	assert.Equal(t, "a", string(text))
}

func TestWrapRunes_WrapsSelectionWithBrackets(t *testing.T) {
	buffer := []rune("hello")

	start, end, text, ok := text.WrappRunes(
		key.Key{Rune: '('},
		1,
		5,
		buffer,
	)

	assert.True(t, ok)
	assert.Equal(t, uint(0), start)
	assert.Equal(t, uint(5), end)
	assert.Equal(t, "(hello)", string(text))
}

func TestWrapRunes_DoesNothingIfRuneIsNotWrapper(t *testing.T) {
	buffer := []rune("hello")

	_, _, text, ok := text.WrappRunes(
		key.Key{Rune: 'a'},
		1,
		4,
		buffer,
	)

	assert.False(t, ok)
	assert.Equal(t, "a", string(text))
}

func TestTextHelper_Apply_UsesFirstMatchingHelper(t *testing.T) {
	helper := text.NewTextTransformer(
		text.AppendSpaceAfterRune,
		text.WrappRunes,
	)

	_, _, text := helper.Apply(
		key.Key{Rune: ','},
		2,
		2,
		[]rune("abc"),
	)

	assert.Equal(t, ", ", string(text))
}

func TestTextHelper_Apply_FallsBackToRuneInsertion(t *testing.T) {
	helper := text.NewTextTransformer(
		text.AppendSpaceAfterRune,
		text.WrappRunes,
	)

	_, _, text := helper.Apply(
		key.Key{Rune: 'x'},
		0,
		0,
		[]rune("abc"),
	)

	assert.Equal(t, "x", string(text))
}
