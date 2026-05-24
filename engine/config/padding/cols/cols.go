package cols

import (
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/marker"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

type FragmentProvider func(winsize.Cols, ...text.Line) text.Fragment

type Option func(*Config)

type Config struct {
	Position style.HorizontalPosition
	Provider FragmentProvider
}

func ResolveConfig(opts ...Option) Config {
	cfg := defaultColsConfig()
	for _, opt := range opts {
		opt(&cfg)
	}
	return cfg
}

func defaultColsConfig() Config {
	return Config{
		Position: style.Left,
		Provider: func(_ winsize.Cols, _ ...text.Line) text.Fragment {
			return *text.NewFragment(marker.DefaultPaddingText)
		},
	}
}

func WithPosition(position style.HorizontalPosition) Option {
	return func(cfg *Config) {
		cfg.Position = position
	}
}

func WithFragment(frag text.Fragment) Option {
	return func(cfg *Config) {
		cfg.Provider = func(_ winsize.Cols, _ ...text.Line) text.Fragment {
			return frag
		}
	}
}

