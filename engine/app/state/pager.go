package state

import "github.com/Rafael24595/go-terminal/engine/helper/math"

type PagerContext struct {
	TargetPage uint
	ActualPage uint
	HasMore    bool
	ForceShow  bool
}

func (s *PagerContext) DecTarget() *PagerContext {
	s.TargetPage = math.SubClampZero(s.TargetPage, 1)
	return s
}

func (s *PagerContext) IncTarget() *PagerContext {
	s.TargetPage += 1
	return s
}

func (s *PagerContext) ConfirmPage(page ...uint) *PagerContext {
	if len(page) > 0 {
		s.TargetPage = page[0]
	}

	s.ActualPage = s.TargetPage
	return s
}
