package pager_test

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/draw"
	"github.com/Rafael24595/go-reacterm-core/engine/app/pager"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
)

type MockStrategy struct {
	EngineCall    uint
	EngineCode    pager.EngineCode
	EngineFunc    pager.EngineFunc
	PredicateCall uint
	PredicateCode pager.PredicateCode
	PredicateBool bool
	PredicateFunc pager.PredicateFunc
}

func (s *MockStrategy) ToStrategy() pager.PagerStrategy {
	return pager.PagerStrategy{
		Engine: pager.Engine{
			Code: s.EngineCode,
			Func: func(dc *draw.DrawContext, ds *draw.DrawState) *draw.DrawState {
				s.EngineCall += 1
				if s.EngineFunc != nil {
					return s.EngineFunc(dc, ds)
				}
				return ds
			},
		},
		Predicate: pager.Predicate{
			Code: s.PredicateCode,
			Func: func(u state.UIState, pc pager.PredicateContext) bool {
				s.PredicateCall += 1
				if s.PredicateFunc != nil {
					return s.PredicateFunc(u, pc)
				}
				return s.PredicateBool
			},
		},
	}
}
