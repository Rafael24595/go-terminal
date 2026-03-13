package action

import "github.com/Rafael24595/go-terminal/engine/core/drawable"

type ActionKind uint8

const (
	ActionMapEach ActionKind = iota
	ActionMapGroup
)

type Focus uint8

const (
	FocusNone   Focus = 0
	FocusHeader Focus = 1 << iota
	FocusBody
	FocusFooter
)

func MergeFocus(styles ...Focus) Focus {
	var merged Focus
	for _, style := range styles {
		merged |= style
	}
	return merged
}

func (s Focus) HasAny(focus ...Focus) bool {
	for _, style := range focus {
		if s&style != 0 {
			return true
		}
	}
	return false
}

func (s Focus) HasNone(focus ...Focus) bool {
	return !s.HasAny(focus...)
}

type ActionFunc func(...drawable.Drawable) []drawable.Drawable

type Action struct {
	Kind     ActionKind
	Focus    Focus
	function ActionFunc
}

func NewAction(kind ActionKind, focus Focus, function ActionFunc) Action {
	return Action{
		Kind:     kind,
		Focus:    focus,
		function: function,
	}
}

func ApplyAction(action Action, items ...drawable.Drawable) []drawable.Drawable {
	if action.Kind == ActionMapGroup {
		return action.function(items...)
	}

	newItems := make([]drawable.Drawable, 0)
	for _, i := range items {
		newItems = append(newItems,
			action.function(i)...,
		)
	}

	return newItems
}
