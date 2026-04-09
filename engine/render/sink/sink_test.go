package sink

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	"github.com/Rafael24595/go-terminal/engine/render/style"
	"github.com/Rafael24595/go-terminal/engine/render/text"
)

func TestApplySinks_PaddingLeft(t *testing.T) {
	spec := style.SpecPaddingLeft(5, "-")

	line := text.EmptyLine()
	line.AddSpec(spec)

	assert.Len(t, 1, line.Text)

	ApplySinks(&line, 80)

	assert.False(t, line.Spec.Kind().HasAny(style.SpcKindPaddingLeft))
	assert.Len(t, 2, line.Text)

	firstFrag := line.Text[0]
	assert.True(t, firstFrag.Spec.Kind().HasAny(style.SpcKindPaddingLeft))
	assert.Equal(t, 5, firstFrag.Spec.Args()[style.KeyPaddingLeftSize].Intd(0))
}
