package node

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
)

func IsKeyRequired(def screen.Definition, ky key.Key) bool {
	return IsActionRequired(def, ky.Code)
}

func IsActionRequired(def screen.Definition, ky key.KeyAction) bool {
	exists := def.RequireKeys.Exists(key.ActionAll)
	if exists {
		return true
	}

	return def.RequireKeys.Exists(ky)
}

func FilterKeyRequired(def screen.Definition, kys ...key.KeyAction) []key.KeyAction {
	filtered := make([]key.KeyAction, 0)
	for _, k := range kys {
		if !IsActionRequired(def, k) {
			filtered = append(filtered, k)
		}
	}
	return filtered
}
