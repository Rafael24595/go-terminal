package key

type ModMask uint8

const (
	ModNone  ModMask = 0
	ModShift ModMask = 1 << iota
	ModAlt
	ModCtrl
)

func MergeMods(mods ...ModMask) ModMask {
	var merged ModMask
	for _, mod := range mods {
		merged |= mod
	}
	return merged
}

func (m ModMask) HasAny(mods ...ModMask) bool {
	for _, mod := range mods {
		if m&mod != 0 {
			return true
		}
	}
	return false
}

func (m ModMask) HasNone(mods ...ModMask) bool {
	return !m.HasAny(mods...)
}
