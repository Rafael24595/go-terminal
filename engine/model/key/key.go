package key

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"
	
	"github.com/Rafael24595/go-terminal/engine/model/help"
)

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

	CustomActionHelp
	CustomActionBack

	CustomActionUndo
	CustomActionRedo

	CustomActionCut
	CustomActionCopy
	CustomActionPaste

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
	'b': NewKeyCode(CustomActionBack, ModAlt),
	'h': NewKeyCode(CustomActionHelp, ModAlt),
	'x': NewKeyCode(CustomActionCut, ModAlt),
	'c': NewKeyCode(CustomActionCopy, ModAlt),
	'v': NewKeyCode(CustomActionPaste, ModAlt),
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

var actionHelpMap = map[KeyAction]help.HelpField{
	ActionArrowUp:    {Code: []string{"↑"}, Detail: "Move up"},
	ActionArrowDown:  {Code: []string{"↓"}, Detail: "Move down"},
	ActionArrowLeft:  {Code: []string{"←"}, Detail: "Move left"},
	ActionArrowRight: {Code: []string{"→"}, Detail: "Move right"},
	ActionHome:       {Code: []string{"HOME", "^A"}, Detail: "Line start"},
	ActionEnd:        {Code: []string{"END", "^E"}, Detail: "Line end"},

	ActionEnter: {Code: []string{"RET"}, Detail: "New line/Accept"},
	ActionTab:   {Code: []string{"TAB"}, Detail: "Next field"},
	ActionEsc:   {Code: []string{"ESC"}, Detail: "Back/Cancel"},
	ActionExit:  {Code: []string{"^C"}, Detail: "Exit"},

	ActionBackspace:      {Code: []string{"BS"}, Detail: "Delete char"},
	ActionDelete:         {Code: []string{"DEL"}, Detail: "Delete forward"},
	ActionDeleteBackward: {Code: []string{"^W"}, Detail: "Delete word"},
	ActionDeleteForward:  {Code: []string{"^D"}, Detail: "Delete word fwd"},

	CustomActionUndo:  {Code: []string{"^G"}, Detail: "Undo"},
	CustomActionRedo:  {Code: []string{"^T"}, Detail: "Redo"},
	CustomActionHelp:  {Code: []string{"M-h"}, Detail: "Help"},
	CustomActionBack:  {Code: []string{"M-b"}, Detail: "Back"},
	CustomActionCut:   {Code: []string{"M-x"}, Detail: "Cut"},
	CustomActionCopy:  {Code: []string{"M-c"}, Detail: "Copy"},
	CustomActionPaste: {Code: []string{"M-v"}, Detail: "Paste"},
}

func ActionsToHelp(actions ...KeyAction) []help.HelpField {
	return ActionsToHelpWithOverride(nil, actions...)
}

func ActionToHelp(action KeyAction) help.HelpField {
	return ActionToHelpWithOverride(nil, action)
}

func ActionsToHelpWithOverride(overrides map[KeyAction]help.HelpField, actions ...KeyAction) []help.HelpField {
	help := make([]help.HelpField, len(actions))
	for i := range actions {
		help[i] = ActionToHelpWithOverride(overrides, actions[i])
	}
	return help
}

func ActionToHelpWithOverride(overrides map[KeyAction]help.HelpField, action KeyAction) help.HelpField {
	if overrides != nil {
		if field, exists := overrides[action]; exists {
			return field
		}
	}

	if str, exist := actionHelpMap[action]; exist {
		return str
	}

	assert.Unreachable("unhandled action: %d", action)

	return help.HelpField{
		Code:   []string{"???"},
		Detail: "Unknown action",
	}
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
