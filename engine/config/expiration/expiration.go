package expiration

import "github.com/Rafael24595/go-reacterm-core/engine/app/screen"

type Expiration struct {
	node *screen.Node
	name string
}

func Persistent() Expiration {
	return Expiration{}
}

func OnNode(node *screen.Node) Expiration {
	return Expiration{
		node: node,
		name: "",
	}
}

func OnName(name string) Expiration {
	return Expiration{
		node: nil,
		name: name,
	}
}

func (e Expiration) On(node *screen.Node) bool {
	if e.node != nil {
		return e.node != node
	}

	if e.name != "" {
		return e.name == node.Name
	}

	return false
}
