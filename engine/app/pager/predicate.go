package pager

import "github.com/Rafael24595/go-terminal/engine/app/state"

type PredicateCode uint16

const (
	CodePredicatePage PredicateCode = iota
	CodePredicateFocus
)

type PredicateContext struct {
	Page     uint
	HasFocus bool
}

type PredicateFunc func(state state.UIState, ctx PredicateContext) bool

type Predicate struct {
	Code PredicateCode
	Func PredicateFunc
}

type PagerStrategy struct {
	Engine    Engine
	Predicate Predicate
}

func PredicatePage() Predicate {
	return Predicate{
		Code: CodePredicatePage,
		Func: func(state state.UIState, ctx PredicateContext) bool {
			return ctx.Page == state.Pager.Page
		},
	}
}

func PredicateFocus() Predicate {
	return Predicate{
		Code: CodePredicateFocus,
		Func: func(_ state.UIState, ctx PredicateContext) bool {
			return ctx.HasFocus
		},
	}
}
