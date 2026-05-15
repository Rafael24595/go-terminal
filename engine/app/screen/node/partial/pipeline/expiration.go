package pipeline

import "github.com/Rafael24595/go-reacterm-core/engine/app/screen"

type expiration struct {
	node *screen.Node
	name string
}

func persistent() expiration {
	return expiration{}
}

func onNode(node *screen.Node) expiration {
	return expiration{
		node: node,
		name: "",
	}
}

func onName(name string) expiration {
	return expiration{
		node: nil,
		name: name,
	}
}

func (e expiration) on(node *screen.Node) bool {
	if e.node != nil {
		return e.node != node
	}

	if e.name != "" {
		return e.name == node.Screen.Name
	}

	return false
}
