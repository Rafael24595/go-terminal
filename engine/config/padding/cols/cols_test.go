package cols

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	"github.com/Rafael24595/go-reacterm-core/engine/render/marker"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

func TestResolveConfigDefaults(t *testing.T) {
	cfg := ResolveConfig()
	assert.Equal(t, style.Left, cfg.Position)

	frag := cfg.Provider(10)
	assert.Equal(t, marker.DefaultPaddingText, frag.Text)
}

func TestWithPosition(t *testing.T) {
	cfg := ResolveConfig(
		WithPosition(style.Right),
	)

	assert.Equal(t, style.Right, cfg.Position)
}

func TestWithFragment(t *testing.T) {
	expected := *text.NewFragment(".")

	cfg := ResolveConfig(
		WithFragment(expected),
	)

	got := cfg.Provider(10)

	assert.Equal(t, expected.Text, got.Text)
}
