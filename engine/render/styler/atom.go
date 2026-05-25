package styler

import (
	"strings"

	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/dict"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
)

type AtomStyler func(string) string

func pa(k style.Atom, s AtomStyler) dict.Pair[style.Atom, AtomStyler] {
	return dict.NewPair(k, s)
}

var Atoms = dict.NewInmutableLinkedMap(
	pa(style.AtmLower, func(text string) string {
		return strings.ToLower(text)
	}),
	pa(style.AtmUpper, func(text string) string {
		return strings.ToUpper(text)
	}),
	pa(style.AtmBold, func(text string) string {
		return text
	}),
	pa(style.AtmSelect, func(text string) string {
		return text
	}),
)

type Atom struct {
	table *dict.LinkedMap[style.Atom, AtomStyler]
}

func NewAtom() *Atom {
	instance := &Atom{}
	return instance.lazyInit()
}

func NewDefaultAtom() *Atom {
	return &Atom{
		table: Atoms.Clone(false),
	}
}

func (a *Atom) lazyInit() *Atom {
	if a.table != nil {
		return a
	}

	a.table = dict.NewLinkedMap[style.Atom, AtomStyler]()
	return a
}

func (a *Atom) Push(pair ...dict.Pair[style.Atom, AtomStyler]) *Atom {
	a.lazyInit()

	a.table.SetPairs(pair...)
	return a
}

func (a *Atom) Apply(text string, styles ...style.Atom) string {
	a.lazyInit()

	merged := style.MergeAtom(styles...)

	for k, p := range a.table.All() {
		if merged.HasAny(k) {
			text = p(text)
		}
	}

	return text
}
