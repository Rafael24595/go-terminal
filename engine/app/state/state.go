package state

type UIState struct {
	Helper HelperContext
	Pager  PagerState
	Stack  *StackContext
}

func NewUIState() *UIState {
	return &UIState{
		Helper: HelperContext{},
		Pager:  PagerState{},
		Stack:  newStackContext(),
	}
}
