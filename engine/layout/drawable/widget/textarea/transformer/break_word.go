package transformer

import (
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

func BreakWord(frags []text.Fragment) []text.Fragment {
	for i := range frags {
		frags[i].AddAtom(style.AtmBreak)
	}
	return frags
}
