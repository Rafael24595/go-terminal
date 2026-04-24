package key

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/model/ascii"
	"github.com/Rafael24595/go-reacterm-core/engine/model/help"
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
	ascii.CTRL_A:     NewKeyCode(ActionHome),
	ascii.CTRL_E:     NewKeyCode(ActionEnd),
	ascii.CTRL_C:     NewKeyCode(ActionExit),
	ascii.CTRL_W:     NewKeyCode(ActionDeleteBackward),
	ascii.CTRL_D:     NewKeyCode(ActionDeleteForward),
	ascii.TAB:        NewKeyCode(ActionTab),
	ascii.ENTER_LF:   NewKeyCode(ActionEnter),
	ascii.ENTER_CR:   NewKeyCode(ActionEnter),
	ascii.DEL:        NewKeyCode(ActionBackspace),
	ascii.BACK_SPACE: NewKeyCode(ActionBackspace),
}

var AltKeyMap = map[rune]*Key{
	'd': NewKeyCode(ActionDeleteForward, ModAlt),
	'b': NewKeyCode(CustomActionBack, ModAlt),
	'h': NewKeyCode(CustomActionHelp, ModAlt),
	'x': NewKeyCode(CustomActionCut, ModAlt),
	'c': NewKeyCode(CustomActionCopy, ModAlt),
	'v': NewKeyCode(CustomActionPaste, ModAlt),
	'z': NewKeyCode(CustomActionUndo, ModAlt),
	'y': NewKeyCode(CustomActionRedo, ModAlt),
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

	ActionRune: {Code: []string{"Text"}, Detail: "Text"},
}

func ActionsToHelp(actions ...KeyAction) []help.HelpField {
	return ActionsToHelpWithOverride(nil, actions...)
}

func ActionToHelp(action KeyAction) *help.HelpField {
	return ActionToHelpWithOverride(nil, action)
}

func ActionsToHelpWithOverride(overrides map[KeyAction]help.HelpField, actions ...KeyAction) []help.HelpField {
	help := make([]help.HelpField, len(actions))
	for i := range actions {
		if action := ActionToHelpWithOverride(overrides, actions[i]); action != nil {
			help[i] = *action
		}
	}
	return help
}

func ActionToHelpWithOverride(overrides map[KeyAction]help.HelpField, action KeyAction) *help.HelpField {
	if action == ActionAll {
		return nil
	}

	if overrides != nil {
		if field, exists := overrides[action]; exists {
			return &field
		}
	}

	if str, exist := actionHelpMap[action]; exist {
		return &str
	}

	assert.Unreachable("unhandled action: %d", action)

	return &help.HelpField{
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
