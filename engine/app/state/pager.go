package state

import (
	"github.com/Rafael24595/go-reacterm-core/engine/helper/math"
)

type PagerContext struct {
	Syncronyzed bool
	modificated bool
	TargetPage  uint
	ActualPage  uint
	HasMore     bool
	ForceShow   bool
}

func (s *PagerContext) DecTarget() *PagerContext {
	s.Syncronyzed = false
	s.modificated = true
	s.TargetPage = math.SubClampZero(s.TargetPage, 1)
	return s
}

func (s *PagerContext) IncTarget() *PagerContext {
	s.Syncronyzed = false
	s.modificated = true
	s.TargetPage += 1
	return s
}

func (s *PagerContext) ConfirmPage(page ...uint) *PagerContext {
	if len(page) > 0 {
		s.TargetPage = page[0]
	}

	if s.modificated &&
		s.TargetPage == s.ActualPage {
		s.Syncronyzed = true
	}

	s.ActualPage = s.TargetPage
	s.modificated = false
	return s
}
