package wrapper_layout

import (
	"github.com/Rafael24595/go-terminal/engine/app/draw"
	"github.com/Rafael24595/go-terminal/engine/app/pager"
	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/app/viewmodel"
	"github.com/Rafael24595/go-terminal/engine/helper/math"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable/primitive/line"
	"github.com/Rafael24595/go-terminal/engine/model/winsize"
	"github.com/Rafael24595/go-terminal/engine/render/text"
)

// TODO: Implement tokenize lines method to prevent line feed injection.
func TerminalApply(state *state.UIState, vm viewmodel.ViewModel, size winsize.Winsize) []text.Line {
	header, footer := vm.InitStaticLayers()

	headerLines := drawStaticLines(header, size)
	footerLines := drawStaticLines(footer, size)

	inputLines := make([]text.Line, 0)
	if input, ok := vm.InitInputLine(size); ok {
		inputLines = drawStaticLines(input, size)
	}

	helperLines := make([]text.Line, 0)
	if helper, ok := vm.InitHelper(size); ok {
		helperLines = drawStaticLines(helper, size)
	}

	static := winsize.Rows(
		len(headerLines) + len(footerLines) + len(inputLines) + len(helperLines),
	)

	if static > size.Rows {
		return []text.Line{
			*text.NewLine("Too low resolution"),
		}
	}

	rest := math.SubClampZero(size.Rows, static)
	remSize := winsize.NewWinsize(rest, size.Cols)
	lines := vm.InitDynamicLayers(remSize)

	dynamicSize := winsize.NewWinsize(rest, size.Cols)
	drawCtx := draw.NewDrawContext(state, dynamicSize)
	drawStt := drawDynamicLines(drawCtx, vm.Pager, lines)

	state.Pager.ConfirmPage(drawStt.Page)
	state.Pager.HasMore = showPagination(drawStt)

	allLines := headerLines
	allLines = append(allLines, drawStt.Buffer...)
	allLines = append(allLines, footerLines...)
	allLines = append(allLines, inputLines...)
	allLines = append(allLines, helperLines...)

	return allLines
}

func drawStaticLines(drawable drawable.Drawable, size winsize.Winsize) []text.Line {
	rows := int(size.Rows)
	cols := int(size.Cols)

	buffer := make([]text.Line, 0)

	content := true
	for content {
		lines, status := drawable.Draw(size)
		content = status

		if len(lines) == 0 {
			break
		}

		for _, lin := range lines {
			buffer = append(buffer,
				line.WrapLineWords(cols, &lin)...,
			)

			if len(buffer) >= rows {
				return buffer
			}
		}
	}

	return buffer
}

func drawDynamicLines(ctx *draw.DrawContext, pager pager.PagerStrategy, drawable drawable.Drawable) *draw.DrawState {
	state := draw.NewDrawStatus(ctx)
	if ctx.Size.Rows == 0 {
		return state
	}

	cols := int(ctx.Size.Cols)

	var rendered []text.Line
	hasNext := true

	for hasNext {
		rendered, hasNext = drawable.Draw(ctx.Size)
		state.HasNext = hasNext

		renderedSize := len(rendered)
		if len(rendered) == 0 {
			continue
		}

		state.Work.Reset()
		state.Work.Add(renderedSize)

		for l, ln := range rendered {
			fixed := line.WrapLineWords(cols, &ln)

			state.Work.Advance()
			state.Work.Add(len(fixed))

			for f, fx := range fixed {
				state.Buffer[state.Cursor] = fx

				state.Work.Advance()

				if f := text.HasFocus(&fx); f {
					state.Focus = f
				}

				state.Cursor += 1
				if winsize.Rows(state.Cursor) < ctx.Size.Rows {
					continue
				}

				if shouldStop(ctx, pager, state) {
					return state
				}

				isLastFixed := f == len(fixed)-1
				isLastRendered := l == len(rendered)-1
				if !isLastFixed || !isLastRendered || hasNext {
					state = pager.Engine.Func(ctx, state)
				}
			}
		}
	}

	return state
}

func shouldStop(ctx *draw.DrawContext, pgr pager.PagerStrategy, stt *draw.DrawState) bool {
	return pgr.Predicate.Func(*ctx.State, pager.PredicateContext{
		Page:     stt.Page,
		HasFocus: stt.Focus,
	})
}

func showPagination(stt *draw.DrawState) bool {
	return stt.Page != 0 || stt.HasNext || (stt.Work.HasWorks() && !stt.Work.Finished())
}
