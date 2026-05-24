package rows

import (
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/marker"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

type FragmentProvider func(winsize.Winsize, ...text.Line) text.Fragment

type Option func(*Config)

type Config struct {
	Position style.VerticalPosition
	Provider FragmentProvider
}

func ResolveConfig(opts ...Option) Config {
	cfg := defaultConfig()
	for _, opt := range opts {
		opt(&cfg)
	}
	return cfg
}

func defaultConfig() Config {
	return Config{
		Position: style.Top,
		Provider: func(_ winsize.Winsize, _ ...text.Line) text.Fragment {
			return *text.EmptyFragment()
		},
	}
}

func WithPosition(position style.VerticalPosition) Option {
	return func(cfg *Config) {
		cfg.Position = position
	}
}

func WithFragment(frag text.Fragment) Option {
	return func(cfg *Config) {
		cfg.Provider = func(_ winsize.Winsize, _ ...text.Line) text.Fragment {
			return frag
		}
	}
}

func WithFillFragment(frag ...string) Option {
	data := marker.DefaultPaddingText
	if len(frag) > 0 {
		data = frag[0]
	}

	return func(cfg *Config) {
		cfg.Provider = func(size winsize.Winsize, lines ...text.Line) text.Fragment {
			measure := text.MaxLineMeasure(size.Cols, lines...)
			return *text.NewFragment(data).
				AddSpec(style.SpecRepeatRight(measure))
		}
	}
}
