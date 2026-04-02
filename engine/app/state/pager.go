package state

import "github.com/Rafael24595/go-terminal/engine/helper/math"

type PagerState struct {
	TargetPage uint
	ActualPage uint
	HasMore    bool
	ForceShow  bool
}

func (s *PagerState) DecTarget() *PagerState {
	s.TargetPage = math.SubClampZero(s.TargetPage, 1)
	return s
}

func (s *PagerState) IncTarget() *PagerState {
	s.TargetPage += 1
	return s
}

func (s *PagerState) ConfirmPage(page ...uint) *PagerState {
	if len(page) > 0 {
		s.TargetPage = page[0]
	}

	s.ActualPage = s.TargetPage
	return s
}
