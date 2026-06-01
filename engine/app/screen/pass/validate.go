package pass

import (
	"errors"
	"fmt"

	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/set"
)

const (
	err_name   = "screen: name is required"
	errf_keys  = "screen %q: Keys is nil"
	errf_tick  = "screen %q: Tick is nil"
	errf_view  = "screen %q: View is nil"
	errf_cycle = "screen %q: Cycle detected"
)

func ValidateStructure() screen.Pass {
	return func(node screen.Node) (screen.Node, error) {
		visited := set.NewSet[string]()

		pending := []screen.Node{node}
		cursor := 0

		for cursor < len(pending) {
			focus := pending[cursor]
			visited.Add(focus.Id())

			if focus.Name == "" {
				return node, errors.New(err_name)
			}

			if focus.Screen.Keys == nil {
				return node, fmt.Errorf(errf_keys, focus.Name)
			}

			if focus.Screen.Tick == nil {
				return node, fmt.Errorf(errf_tick, focus.Name)
			}

			if focus.Screen.View == nil {
				return node, fmt.Errorf(errf_view, focus.Name)
			}

			children := focus.Children()
			for i := range children {
				child := children[i]
				if visited.Has(child.Id()) {
					return node, fmt.Errorf(errf_cycle, child.Name)
				}

				pending = append(pending, child)
			}

			cursor += 1
		}

		return node, nil
	}
}
