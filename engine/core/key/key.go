package key

const (
	CTRL_C   = 0x03
	CTRL_W   = 0x17
	TAB      = 0x09
	ENTER_LF = '\n'
	ENTER_CR = '\r'
	ESC      = 0x1b
	DEL      = 0x7f
	BS       = 0x08
	TILDE    = '~'
)

const (
	VK_SHIFT = 0x10
)

type KeyCode int

const (
	KeyRune KeyCode = iota
	KeyCtrlC
	KeyDeleteWordBackward
	KeyDeleteWordForward
	KeyTab
	KeyEnter
	KeyBackspace
	KeyArrowUp
	KeyArrowDown
	KeyArrowLeft
	KeyArrowRight
	KeyHome
	KeyEnd
	KeyDelete
	KeyAll
)

type ModMask uint8

const (
	ModNone  ModMask = 0
	ModShift ModMask = 1 << iota
	ModAlt
	ModCtrl
)

func (m ModMask) Has(mod ModMask) bool {
	return m&mod != 0
}

type Key struct {
	Code KeyCode
	Mod  ModMask
	Rune rune
}

func NewKeysCode(codes ...KeyCode) []Key {
	keys := make([]Key, len(codes))

	for i, v := range codes {
		keys[i] = Key{Code: v}
	}

	return keys
}

func NewKeyCode(code KeyCode, mods ...ModMask) *Key {
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
		Code: KeyRune,
		Mod:  ModNone,
		Rune: r,
	}
}

func NewKeySpace() *Key {
	return NewKeyRune(' ')
}
