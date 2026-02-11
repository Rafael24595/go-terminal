package key

const (
	CTRL_A = 0x01
	CTRL_E = 0x05
	CTRL_D = 0x04
	CTRL_C = 0x03
	CTRL_G = 0x07
	CTRL_T = 0x14
	CTRL_W = 0x17
)

const (
	ESC        = 0x1b
	DEL        = 0x7f
	TAB        = 0x09
	ENTER_LF   = 0x0A
	ENTER_CR   = 0x0D
	TILDE      = 0x7E
	BACK_SPACE = 0x08
)

const (
	VK_SHIFT = 0x10
)

type KeyAction int

const (
	ActionRune KeyAction = iota
	
	ActionEsc
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

var ControlKeyMap = map[rune]*Key{
	CTRL_A:     NewKeyCode(ActionHome),
	CTRL_E:     NewKeyCode(ActionEnd),
	CTRL_C:     NewKeyCode(ActionExit),
	CTRL_W:     NewKeyCode(ActionDeleteBackward),
	CTRL_D:     NewKeyCode(ActionDeleteForward),
	CTRL_G:     NewKeyCode(CustomActionUndo),
	CTRL_T:     NewKeyCode(CustomActionRedo),
	TAB:        NewKeyCode(ActionTab),
	ENTER_LF:   NewKeyCode(ActionEnter),
	ENTER_CR:   NewKeyCode(ActionEnter),
	DEL:        NewKeyCode(ActionBackspace),
	BACK_SPACE: NewKeyCode(ActionBackspace),
}

var AltKeyMap = map[rune]*Key{
	'd': NewKeyCode(ActionDeleteForward, ModAlt),
}

var CsiFinalMap = map[rune]KeyAction{
	'A': ActionArrowUp,
	'B': ActionArrowDown,
	'C': ActionArrowRight,
	'D': ActionArrowLeft,
	'H': ActionHome,
	'F': ActionEnd,
}

var CsiTildeMap = map[string]KeyAction{
	"3": ActionDelete,
	"1": ActionHome,
	"7": ActionHome,
	"4": ActionEnd,
	"8": ActionEnd,
}

var actionMap = map[KeyAction][]string{
	ActionEsc:            {"ESC"},
	ActionExit:           {"CTRL+C"},
	ActionDeleteBackward: {"CTRL+W"},
	ActionDeleteForward:  {"CTRL+D"},
	ActionTab:            {"TAB"},
	ActionEnter:          {"ENTER"},
	ActionBackspace:      {"BACKSPACE"},
	ActionArrowUp:        {"UP"},
	ActionArrowDown:      {"DOWN"},
	ActionArrowLeft:      {"LEFT"},
	ActionArrowRight:     {"RIGHT"},
	ActionHome:           {"HOME", "CTRL+A"},
	ActionEnd:            {"END", "CTRL+E"},
	ActionDelete:         {"DELETE"},
	CustomActionUndo:     {"CTRL+G"},
	CustomActionRedo:     {"CTRL+T"},
	ActionAll:            {"$_SYSTEM_ALL"},
}

func ActionToString(action KeyAction) []string {
	if str, exist := actionMap[action]; exist {
		return str
	}
	return []string{"rune"}
}

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
