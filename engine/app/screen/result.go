package screen

import "github.com/Rafael24595/go-reacterm-core/engine/app/state"

type Result struct {
	Isolate bool
	Node    *Node
	Pager   state.PagerContext
}

func ResultFromNode(node *Node) Result {
	return Result{
		Isolate: false,
		Node:    node,
		Pager:   state.PagerContext{},
	}
}

func ResultFromUIState(stt *state.UIState) Result {
	return Result{
		Isolate: false,
		Node:    nil,
		Pager:   stt.Pager,
	}
}

func EmptyResult() Result {
	return Result{
		Isolate: false,
		Node:    nil,
		Pager:   state.PagerContext{},
	}
}
