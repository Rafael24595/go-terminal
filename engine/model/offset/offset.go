package offset

import "github.com/Rafael24595/go-reacterm-core/engine/helper/math"

type Offset uint32

func (r Offset) Clamp(o Offset) Offset {
	return math.SubClampZero(r, o)
}
