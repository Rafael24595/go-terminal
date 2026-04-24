package screen

import (
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
)

func IsKeyRequired(def Definition, ky key.Key) bool {
	for _, v := range def.RequireKeys {
		if v.Code == key.ActionAll || v.Code == ky.Code {
			return true
		}
	}
	return false
}

func IsActionRequired(def Definition, ky key.KeyAction) bool {
	for _, v := range def.RequireKeys {
		if v.Code == key.ActionAll || v.Code == ky {
			return true
		}
	}
	return false
}

func FilterKeyRequired(def Definition, kys ...key.KeyAction) []key.KeyAction {
	filtered := make([]key.KeyAction, 0)
	for _, k := range kys {
		if !IsActionRequired(def, k) {
			filtered = append(filtered, k)
		}
	}
	return filtered
}
