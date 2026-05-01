package line

import "github.com/Rafael24595/go-reacterm-core/engine/helper"

type indexMeta struct {
	sufix      string
	prefixBody string
	digits     uint16
	totalWidth uint32
}

func (i indexMeta) header(index int) string {
	return helper.Right(index, int(i.digits)) + i.sufix
}

func (i indexMeta) body() string {
	return i.prefixBody + i.sufix
}
