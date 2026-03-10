package state

type PagerMode uint16

const (
	PagerModePage PagerMode = iota
	PagerModeCursor
	PagerModeFocus
)

type UIState struct {
	Helper HelperContext
	Pager  PagerContext
}

func NewUIState() *UIState {
	return &UIState{
		Helper: HelperContext{},
		Pager:  PagerContext{},
	}
}

type HelperContext struct {
	ShowHelp bool
}

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
