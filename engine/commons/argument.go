package commons

import (
	"fmt"
	"strconv"
	"strings"
)

type Argument struct {
	item any
}

func ArgumentFrom(item any) *Argument {
	return &Argument{
		item: item,
	}
}

func (a Argument) Bool() (bool, bool) {
	switch v := a.item.(type) {
	case bool:
		return v, true
	case int, int8, int16, int32, int64:
		return a.Int64d(0) != 0, true
	case uint, uint8, uint16, uint32, uint64:
		return a.Int64d(0) != 0, true
	case float32, float64:
		return a.Float64d(0) != 0, true
	case string:
		val, err := strconv.ParseBool(strings.ToLower(v))
		if err == nil {
			return val, true
		}
	}
	return false, false
}

func (a Argument) Boold(def bool) bool {
	if v, ok := a.Bool(); ok {
		return v
	}
	return def
}

func (a Argument) String() string {
	switch v := a.item.(type) {
	case nil:
		return ""
	case string:
		return v
	case bool:
		return strconv.FormatBool(v)
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", v)
	case uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	}
	return fmt.Sprintf("%v", a.item)
}

func (a Argument) Int() (int, bool) {
	if v, ok := a.Int64(); ok {
		return int(v), true
	}
	return 0, false
}

func (a Argument) Intd(def int) int {
	if v, ok := a.Int64(); ok {
		return int(v)
	}
	return def
}

func (a Argument) Int8() (int8, bool) {
	if v, ok := a.Int64(); ok {
		return int8(v), true
	}
	return 0, false
}

func (a Argument) Int8d(def int8) int8 {
	if v, ok := a.Int8(); ok {
		return v
	}
	return def
}

func (a Argument) Int16() (int16, bool) {
	if v, ok := a.Int64(); ok {
		return int16(v), true
	}
	return 0, false
}

func (a Argument) Int16d(def int16) int16 {
	if v, ok := a.Int16(); ok {
		return v
	}
	return def
}

func (a Argument) Int32() (int32, bool) {
	if v, ok := a.Int64(); ok {
		return int32(v), true
	}
	return 0, false
}

func (a Argument) Int32d(def int32) int32 {
	if v, ok := a.Int32(); ok {
		return v
	}
	return def
}

func (a Argument) Int64() (int64, bool) {
	switch v := a.item.(type) {
	case int:
		return int64(v), true
	case int8:
		return int64(v), true
	case int16:
		return int64(v), true
	case int32:
		return int64(v), true
	case int64:
		return v, true
	case uint:
		return int64(v), true
	case uint8:
		return int64(v), true
	case uint16:
		return int64(v), true
	case uint32:
		return int64(v), true
	case uint64:
		return int64(v), true
	case float32:
		return int64(v), true
	case float64:
		return int64(v), true
	case bool:
		if v {
			return 1, true
		}
		return 0, true
	case string:
		val, err := strconv.ParseInt(v, 10, 64)
		if err == nil {
			return val, true
		}
	}
	return 0, false
}

func (a Argument) Int64d(def int64) int64 {
	if v, ok := a.Int64(); ok {
		return v
	}
	return def
}

func (a Argument) Uint() (uint, bool) {
	if v, ok := a.Int64(); ok && v >= 0 {
		return uint(v), true
	}
	return 0, false
}

func (a Argument) Uintd(def uint) uint {
	if v, ok := a.Uint(); ok {
		return v
	}
	return def
}

func (a Argument) Uint8() (uint8, bool) {
	if v, ok := a.Int64(); ok && v >= 0 {
		return uint8(v), true
	}
	return 0, false
}

func (a Argument) Uint8d(def uint8) uint8 {
	if v, ok := a.Uint8(); ok {
		return v
	}
	return def
}

func (a Argument) Uint16() (uint16, bool) {
	if v, ok := a.Int64(); ok && v >= 0 {
		return uint16(v), true
	}
	return 0, false
}

func (a Argument) Uint16d(def uint16) uint16 {
	if v, ok := a.Uint16(); ok {
		return v
	}
	return def
}

func (a Argument) Uint32() (uint32, bool) {
	if v, ok := a.Int64(); ok && v >= 0 {
		return uint32(v), true
	}
	return 0, false
}

func (a Argument) Uint32d(def uint32) uint32 {
	if v, ok := a.Uint32(); ok {
		return v
	}
	return def
}

func (a Argument) Uint64() (uint64, bool) {
	if v, ok := a.Int64(); ok && v >= 0 {
		return uint64(v), true
	}
	return 0, false
}

func (a Argument) Uint64d(def uint64) uint64 {
	if v, ok := a.Uint64(); ok {
		return v
	}
	return def
}

func (a Argument) Float32() (float32, bool) {
	if v, ok := a.Float64(); ok {
		return float32(v), true
	}
	return 0, false
}

func (a Argument) Float32d(def float32) float32 {
	if v, ok := a.Float64(); ok {
		return float32(v)
	}
	return def
}

func (a Argument) Float64() (float64, bool) {
	switch v := a.item.(type) {
	case int:
		return float64(v), true
	case int8:
		return float64(v), true
	case int16:
		return float64(v), true
	case int32:
		return float64(v), true
	case int64:
		return float64(v), true
	case uint:
		return float64(v), true
	case uint8:
		return float64(v), true
	case uint16:
		return float64(v), true
	case uint32:
		return float64(v), true
	case uint64:
		return float64(v), true
	case float32:
		return float64(v), true
	case float64:
		return v, true
	case bool:
		if v {
			return 1, true
		}
		return 0, true
	case string:
		val, err := strconv.ParseFloat(v, 64)
		if err == nil {
			return val, true
		}
	}
	return 0, false
}

func (a Argument) Float64d(def float64) float64 {
	if v, ok := a.Float64(); ok {
		return v
	}
	return def
}

func (a Argument) Parse(parse func(string) (any, error)) (any, bool) {
	v, err := parse(a.String())
	if err != nil {
		return nil, false
	}
	return v, true
}
