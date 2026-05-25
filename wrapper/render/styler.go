package wrapper_render

import (
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/dict"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/styler"

	wrapper_ansi "github.com/Rafael24595/go-reacterm-core/wrapper/ansi"
)

func pa(k style.Atom, s styler.AtomStyler) dict.Pair[style.Atom, styler.AtomStyler] {
	return dict.NewPair(k, s)
}

var Atoms = dict.NewInmutableLinkedMap(
	pa(style.AtmBold, func(text string) string {
		if text == "" {
			return text
		}
		return wrapper_ansi.Bold + text + wrapper_ansi.NormalWeight
	}),
	pa(style.AtmSelect, func(text string) string {
		if text == "" {
			return text
		}
		return wrapper_ansi.Reverse + text + wrapper_ansi.NoReverse
	}),
)
