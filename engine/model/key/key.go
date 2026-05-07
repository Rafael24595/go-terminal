package key

type Key struct {
	Code Action
	Mod  ModMask
	Rune rune
}

func NewKeysCode(codes ...Action) []Key {
	keys := make([]Key, len(codes))

	for i, v := range codes {
		keys[i] = Key{Code: v}
	}

	return keys
}

func NewKeyCode(code Action, mods ...ModMask) *Key {
	var mod ModMask
	for _, m := range mods {
		mod |= m
	}

	return &Key{
		Code: code,
		Mod:  mod,
		Rune: ' ',
	}
}

func NewKeyRune(r rune) *Key {
	return &Key{
		Code: ActionRune,
		Mod:  ModNone,
		Rune: r,
	}
}

func NewKeySpace() *Key {
	return NewKeyRune(' ')
}
