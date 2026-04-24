package pager

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/draw"
	"github.com/Rafael24595/go-reacterm-core/engine/helper/math"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

type EngineCode uint16

const (
	CodeEnginePaged EngineCode = iota
	CodeEngineScroll
)

type EngineFunc func(*draw.DrawContext, *draw.DrawState) *draw.DrawState

type Engine struct {
	Code EngineCode
	Func EngineFunc
}

func EnginePage() Engine {
	return Engine{
		Code: CodeEnginePaged,
		Func: func(ctx *draw.DrawContext, stt *draw.DrawState) *draw.DrawState {
			stt.Buffer = make([]text.Line, ctx.Size.Rows)
			stt.Cursor = 0

			stt.Focus = false
			stt.Page += 1

			return stt
		},
	}
}

func EngineScroll() Engine {
	return Engine{
		Code: CodeEngineScroll,
		Func: func(ctx *draw.DrawContext, stt *draw.DrawState) *draw.DrawState {
			if len(stt.Buffer) == 0 {
				return stt
			}

			copy(stt.Buffer, stt.Buffer[1:])
			stt.Buffer[len(stt.Buffer)-1] = text.Line{}
			stt.Cursor = math.SubClampZero(stt.Cursor, 1)

			stt.Focus = false
			stt.Page += 1

			return stt
		},
	}
}
