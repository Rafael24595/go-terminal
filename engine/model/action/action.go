package action

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
