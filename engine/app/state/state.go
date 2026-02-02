package state

type PagerState struct {
	Enabled bool
	Page    uint
}

func NewPageState(page uint) PagerState {
	return PagerState{
		Enabled: true,
		Page:    page,
	}
}

func EmptyPagerState() PagerState {
	return PagerState{
		Enabled: false,
		Page:    0,
	}
}

type CursorState struct {
	Enabled bool
	Cursor  uint
	Offset  uint
}

func NewCursorState(cursor uint) CursorState {
	return CursorState{
		Enabled: true,
		Cursor:  cursor,
		Offset:  0,
	}
}

func EmptyCursorState() CursorState {
	return CursorState{
		Enabled: false,
		Cursor:  0,
		Offset:  0,
	}
}

type UIState struct {
	Pager  PagerState
	Cursor CursorState
}

func NewUIState() *UIState {
	return &UIState{
		Pager:  EmptyPagerState(),
		Cursor: EmptyCursorState(),
	}
}
