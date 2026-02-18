package style

type Atom uint8

const (
	AtmNone Atom = 0
	AtmBold Atom = 1 << iota
	AtmUpper
	AtmLower
	AtmSelect
)

func MergeAtom(styles ...Atom) Atom {
	var merged Atom
	for _, style := range styles {
		merged |= style
	}
	return merged
}

func (s Atom) HasAny(styles ...Atom) bool {
	for _, style := range styles {
		if s&style != 0 {
			return true
		}
	}
	return false
}

func (s Atom) HasNone(styles ...Atom) bool {
	return !s.HasAny(styles...)
}
