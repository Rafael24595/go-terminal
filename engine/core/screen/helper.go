package screen

import (
	"github.com/Rafael24595/go-terminal/engine/core/key"
)

func IsKeyRequired(def Definition, ky key.Key) bool {
	for _, v := range def.RequireKeys {
		if v.Code == key.ActionAll || v.Code == ky.Code {
			return true
		}
	}
	return false
}
