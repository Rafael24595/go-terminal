package state

type PagerMode uint16

const (
	PagerModePage PagerMode = iota
	PagerModeCursor
	PagerModeFocus
)

type PagerContext struct {
	Page     uint
	Cursor   uint
	Focus    bool
	ShowPage bool
	RestData bool
}

type PagerStrategy struct {
	Mode  PagerMode
	Match func(state UIState, ctx PagerContext) bool
}

func NewPagePager() PagerStrategy {
	return PagerStrategy{
		Mode: PagerModePage,
		Match: func(state UIState, ctx PagerContext) bool {
			return ctx.Page == state.Pager.Page
		},
	}
}

// Deprecated
func NewCursorPager(cursor uint) PagerStrategy {
	return PagerStrategy{
		Mode: PagerModeCursor,
		Match: func(stt UIState, ctx PagerContext) bool {
			return ctx.Cursor >= cursor
		},
	}
}

func NewFocusPager() PagerStrategy {
	return PagerStrategy{
		Mode: PagerModeFocus,
		Match: func(_ UIState, ctx PagerContext) bool {
			return ctx.Focus
		},
	}
}
