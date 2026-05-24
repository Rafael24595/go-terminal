package options

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

func TestWithPosition(t *testing.T) {
	cfg := defaultRowsConfig()

	WithPosition(style.Bottom)(&cfg)

	assert.Equal(t, style.Bottom, cfg.Position)
}

func TestWithFragment(t *testing.T) {
	cfg := defaultRowsConfig()

	WithFragment(
		*text.NewFragment("golang"),
	)(&cfg)

	frag := cfg.Provider(
		winsize.New(10, 20),
	)

	assert.Equal(t, "golang", frag.Text)
}

func TestWithFillFragment(t *testing.T) {
	cfg := defaultRowsConfig()

	WithFillFragment(".")(&cfg)

	lines := []text.Line{
		*text.NewLine("Golang"),
	}

	frag := cfg.Provider(
		winsize.New(10, 20),
		lines...,
	)

	assert.Equal(t, ".", frag.Text)
	assert.True(t, frag.Spec.Kind().HasAny(style.SpcKindRepeatRight))
	assert.Equal(t, "6", frag.Spec.Args()[style.KeyRepeatRightSize].Stringf())
}
