package screen

import (
	"github.com/Rafael24595/go-reacterm-core/engine/model/help"
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
)

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
