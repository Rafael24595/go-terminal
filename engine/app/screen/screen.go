package screen

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/set"
)

type DefinitionFunc func() Definition
type UpdateFunc func(*state.UIState, ScreenEvent) Result
type ViewFunc func(state.UIState) viewmodel.ViewModel

type Screen struct {
	Name       string
	meta       ScreenMeta
	Stack      set.Set[string]
	Definition DefinitionFunc
	Update     UpdateFunc
	View       ViewFunc
}

func (s Screen) Compile(middleware ...ScreenPass) (Screen, error) {
	screen := s
	meta := s.meta

	for _, m := range middleware {
		nextScreen, nextMeta, err := m(screen, meta)
		if err != nil {
			return screen, err
		}

		screen = nextScreen
		meta = nextMeta

		screen.meta = meta
	}

	return screen, nil
}
