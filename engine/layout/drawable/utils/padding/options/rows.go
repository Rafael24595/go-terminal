package options

import (
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/marker"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

type FragmentProvider func(winsize.Winsize, ...text.Line) text.Fragment

type RowsOption func(*RowsConfig)

type RowsConfig struct {
	Position style.VerticalPosition
	Provider FragmentProvider
}

func ResolveRowsConfig(opts ...RowsOption) RowsConfig {
	cfg := defaultRowsConfig()
	for _, opt := range opts {
		opt(&cfg)
	}
	return cfg
}

func defaultRowsConfig() RowsConfig {
	return RowsConfig{
		Position: style.Top,
		Provider: func(_ winsize.Winsize, _ ...text.Line) text.Fragment {
			return *text.EmptyFragment()
		},
	}
}

func WithPosition(position style.VerticalPosition) RowsOption {
	return func(cfg *RowsConfig) {
		cfg.Position = position
	}
}

func WithFragment(frag text.Fragment) RowsOption {
	return func(cfg *RowsConfig) {
		cfg.Provider = func(_ winsize.Winsize, _ ...text.Line) text.Fragment {
			return frag
		}
	}
}

func WithFillFragment(frag ...string) RowsOption {
	data := marker.DefaultPaddingText
	if len(frag) > 0 {
		data = frag[0]
	}

	return func(cfg *RowsConfig) {
		cfg.Provider = func(size winsize.Winsize, lines ...text.Line) text.Fragment {
			measure := text.MaxLineMeasure(size.Cols, lines...)
			return *text.NewFragment(data).
				AddSpec(style.SpecRepeatRight(measure))
		}
	}
}
