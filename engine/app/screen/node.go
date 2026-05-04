package screen

import "github.com/Rafael24595/go-reacterm-core/engine/commons/structure/set"

type Node struct {
	id       string
	Screen   Screen
	meta     Meta
	Stack    set.Set[string]
	children []Node
}

func (n Node) Id() string {
	return n.id
}

func (n Node) Meta() Meta {
	return n.meta
}

func (n Node) Children() []Node {
	return n.children
}

func (n Node) Compile(pass ...Pass) (Node, error) {
	screen := n

	for _, m := range pass {
		nextScreen, err := m(screen)
		if err != nil {
			return screen, err
		}

		screen = nextScreen
	}

	return screen, nil
}
