package style

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
)

func TestEraseSpec_DeleteExists(t *testing.T) {
	scp := MergeSpec(
		SpecPaddingLeft(10, " "),
		SpecFill(100),
	)

	modified, removed := EraseSpec(scp, SpcKindPaddingLeft)

	assert.Equal(t, SpcKindFill, modified.kind)
	assert.NotContains(t, modified.args, KeyPaddingLeftSize)
	assert.Contains(t, modified.args, KeyFillSize)

	assert.Equal(t, SpcKindPaddingLeft, removed.kind)
	assert.Equal(t, 10, removed.args[KeyPaddingLeftSize].Intd(0))
	assert.NotContains(t, removed.args, KeyFillSize)
}

func TestEraseSpec_DeleteNonExists(t *testing.T) {
	scp := MergeSpec(
		SpecPaddingLeft(10, " "),
		SpecFill(100),
	)

	modified, removed := EraseSpec(scp, SpcKindTrimLeft)

	assert.Equal(t, scp.kind, modified.kind)
	assert.Equal(t, len(scp.args), len(modified.args))

	assert.Equal(t, SpcKindNone, removed.kind)
	assert.Equal(t, 0, len(removed.args))
}

func TestEraseSpec_DeleteMultiple(t *testing.T) {
	scp := MergeSpec(
		SpecPaddingLeft(10, " "),
		SpecFill(100),
	)

	toRemove := SpcKindPaddingLeft | SpcKindFill | SpcKindRepeatRight
	modified, removed := EraseSpec(scp, toRemove)

	assert.Equal(t, SpcKindNone, modified.kind)
	assert.Equal(t, 0, len(modified.args))

	assert.Equal(t, SpcKindPaddingLeft|SpcKindFill, removed.kind)
	assert.Equal(t, 3, len(removed.args))
}
