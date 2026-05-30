package key

import (
	"github.com/Rafael24595/go-reacterm-core/engine/model/ascii"
)

type Action int

const (
	ActionRune Action = iota

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

	ActionPageUp
	ActionPageDown

	CustomActionHelp
	CustomActionBack

	CustomActionUndo
	CustomActionRedo

	CustomActionCut
	CustomActionCopy
	CustomActionPaste

	CustomActionPointer

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
	'p': NewKeyCode(CustomActionPointer, ModAlt),
}

var CsiFinalMap = map[rune]Action{
	'A': ActionArrowUp,
	'B': ActionArrowDown,
	'C': ActionArrowRight,
	'D': ActionArrowLeft,
	'H': ActionHome,
	'F': ActionEnd,
}

var CsiTildeMap = map[string]Action{
	"3": ActionDelete,
	"1": ActionHome,
	"7": ActionHome,
	"4": ActionEnd,
	"8": ActionEnd,
	"5": ActionPageUp,
	"6": ActionPageDown,
}
