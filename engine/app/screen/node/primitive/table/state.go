package table

import (
	"github.com/Rafael24595/go-reacterm-core/engine/model/param"
)

const ArgTableState param.Typed[State] = "table_state"

type State struct {
	Row uint16
	Col uint16
}
