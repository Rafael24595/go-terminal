package key

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/dict"
)

type Descriptor struct {
	Code   []string
	Detail string
}

var actionHelpMap = map[Action]Descriptor{
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

	CustomActionPointer: {Code: []string{"M-p"}, Detail: "Switch gutter"},

	ActionRune: {Code: []string{"Text"}, Detail: "Text"},
}

func ResolveDescriptors(actions ...Action) *dict.LinkedMap[Action, Descriptor] {
	return ResolveDescriptorsWithDefaults(nil, actions...)
}

func ResolveDescriptorsWithDefaults(
	defaults map[Action]Descriptor,
	actions ...Action,
) *dict.LinkedMap[Action, Descriptor] {
	help := dict.NewLinkedMap[Action, Descriptor]()
	for _, a := range actions {
		if action := resolveDescriptor(defaults, a); action != nil {
			help.Set(a, *action)
		}
	}

	return help
}

func resolveDescriptor(
	defaults map[Action]Descriptor,
	action Action,
) *Descriptor {
	if action == ActionAll {
		return nil
	}

	if defaults != nil {
		if field, exists := defaults[action]; exists {
			return &field
		}
	}

	if str, exist := actionHelpMap[action]; exist {
		return &str
	}

	assert.Unreachable("unhandled action: %d", action)

	return &Descriptor{
		Code:   []string{"???"},
		Detail: "Unknown action",
	}
}
