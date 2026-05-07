package screen

import (
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/dict"
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
)

type Definition struct {
	RequireKeys *dict.LinkedMap[key.Action, key.Key]
	Descriptor  *dict.LinkedMap[key.Action, key.Descriptor]
}

func (d Definition) Merge(other Definition) Definition {
	required := d.RequireKeys.Clone()
	required.Merge(other.RequireKeys)

	descriptor := other.Descriptor.Clone()
	descriptor.Supplement(d.Descriptor)

	return Definition{
		RequireKeys: required,
		Descriptor:  descriptor,
	}
}

func NewDefinition(
	defaults map[key.Action]key.Descriptor,
	actions []key.Action,
) Definition {
	keys := key.NewKeysCode(actions...)

	required := dict.NewLinkedMap[key.Action, key.Key]()
	for _, v := range keys {
		required.Set(v.Code, v)
	}

	descriptor := key.ResolveDescriptorsWithDefaults(
		defaults, actions...,
	)

	return Definition{
		RequireKeys: required,
		Descriptor:  descriptor,
	}
}

func EmptyDefinition() Definition {
	return Definition{
		RequireKeys: dict.NewLinkedMap[key.Action, key.Key](),
		Descriptor:  dict.NewLinkedMap[key.Action, key.Descriptor](),
	}
}

func DefinitionFromActions(actions ...key.Action) Definition {
	return NewDefinition(make(map[key.Action]key.Descriptor), actions)
}
