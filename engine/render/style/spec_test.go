package style

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	"github.com/Rafael24595/go-reacterm-core/engine/commons"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
)

func TestEraseSpec_DeleteExists(t *testing.T) {
	scp := MergeSpec(
		SpecPaddingLeft(10, " "),
		SpecFill(100),
	)

	modified, removed := EraseSpec(scp, SpcKindPaddingLeft)
	size := commons.Mapd[winsize.Cols](removed.args[KeyPaddingLeftSize], 0)

	assert.Equal(t, SpcKindFill, modified.kind)
	assert.NotContains(t, modified.args, KeyPaddingLeftSize)
	assert.Contains(t, modified.args, KeyFillSize)

	assert.Equal(t, SpcKindPaddingLeft, removed.kind)
	assert.Equal(t, 10, size)
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
