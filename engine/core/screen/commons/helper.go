package commons

import (
	"github.com/Rafael24595/go-terminal/engine/core/key"
	"github.com/Rafael24595/go-terminal/engine/core/screen"
)

func isKeyRequired(def screen.Definition, ky key.Key) bool {
	for _, v := range def.RequireKeys {
		if v.Code == key.ActionAll || v.Code == ky.Code {
			return true
		}
	}
	return false
}
