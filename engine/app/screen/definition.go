package screen

import (
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/dict"
	"github.com/Rafael24595/go-reacterm-core/engine/model/help"
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
)

type Definition struct {
	RequireKeys *dict.LinkedMap[key.KeyAction, key.Key]
}

func (d Definition) Merge(other Definition) Definition {
    required := d.RequireKeys.Clone()
    required.Merge(other.RequireKeys)
	
    return Definition{
        RequireKeys: required,
    }
}

func DefinitionFromKeys(keys ...key.Key) Definition {
	required := dict.NewLinkedMap[key.KeyAction, key.Key]()
	for _, v := range keys {
		required.Set(v.Code, v)
	}

	return Definition{
		RequireKeys: required,
	}
}

type DefinitionSources struct {
	Overrides  map[key.KeyAction]help.HelpField
	Actions    []key.KeyAction
	Keys       []key.Key
	Definition Definition
}

func NewDefinitionSources(
	overrides map[key.KeyAction]help.HelpField,
	actions []key.KeyAction,
) DefinitionSources {
	keys := key.NewKeysCode(actions...)
	definition := DefinitionFromKeys(keys...)
	return DefinitionSources{
		Overrides:  overrides,
		Actions:    actions,
		Keys:       keys,
		Definition: definition,
	}
}
