package state

type UIState struct {
	Helper HelperContext
	Pager  PagerContext
	Stack  *StackContext
}

func NewUIState() *UIState {
	return &UIState{
		Helper: HelperContext{},
		Pager:  PagerContext{},
		Stack:  newStackContext(),
	}
}
