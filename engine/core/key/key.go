package key

const (
	CTRL_C   = 0x03
	CTRL_W   = 0x17
	CTRL_G   = 0x07
	CTRL_T   = 0x14
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

type KeyAction int

const (
	ActionRune KeyAction = iota
	ActionExit
	ActionDeleteBackward
	ActionDeleteForward
	ActionTab
	ActionEnter
	ActionBackspace
	ActionArrowUp
	ActionArrowDown
	ActionArrowLeft
	ActionArrowRight
	ActionHome
	ActionEnd
	ActionDelete
	CustomActionUndo
	CustomActionRedo
	ActionAll
)

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

type Key struct {
	Code KeyAction
	Mod  ModMask
	Rune rune
}

func NewKeysCode(codes ...KeyAction) []Key {
	keys := make([]Key, len(codes))

	for i, v := range codes {
		keys[i] = Key{Code: v}
	}

	return keys
}

func NewKeyCode(code KeyAction, mods ...ModMask) *Key {
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
