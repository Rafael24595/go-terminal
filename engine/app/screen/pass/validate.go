package pass

import (
	"errors"
	"fmt"

	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/set"
)

const (
	err_name        = "screen: name is required"
	errf_definition = "screen %q: Definition is nil"
	errf_update     = "screen %q: Update is nil"
	errf_view       = "screen %q: View is nil"
	errf_cycle      = "screen %q: Cycle detected"
)

func ValidateStructure() screen.Pass {
	return func(node screen.Node) (screen.Node, error) {
		visited := set.NewSet[string]()

		pending := []screen.Node{node}
		cursor := 0

		for cursor < len(pending) {
			focus := pending[cursor]
			visited.Add(focus.Id())

			if focus.Screen.Name == "" {
				return node, errors.New(err_name)
			}

			if focus.Screen.Definition == nil {
				return node, fmt.Errorf(errf_definition, focus.Screen.Name)
			}

			if focus.Screen.Update == nil {
				return node, fmt.Errorf(errf_update, focus.Screen.Name)
			}

			if focus.Screen.View == nil {
				return node, fmt.Errorf(errf_view, focus.Screen.Name)
			}

			children := focus.Children()
			for i := range children {
				child := children[i]
				if visited.Has(child.Id()) {
					return node, fmt.Errorf(errf_cycle, child.Screen.Name)
				}

				pending = append(pending, child)
			}

			cursor += 1
		}

		return node, nil
	}
}
