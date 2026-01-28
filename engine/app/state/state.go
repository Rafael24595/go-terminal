package state

type LayoutState struct {
	Page       uint
	Pagination bool
}

func newLayoutState() LayoutState {
	return LayoutState{
		Page:       0,
		Pagination: false,
	}
}

type InteractionState struct {
	Cursor uint
	Offset uint
}

func newInteractionState() InteractionState {
	return InteractionState{
		Cursor: 0,
		Offset: 0,
	}
}

type UIState struct {
	Layout      LayoutState
	Interaction InteractionState
}

func NewUIState() *UIState {
	return &UIState{
		Layout:      newLayoutState(),
		Interaction: newInteractionState(),
	}
}
