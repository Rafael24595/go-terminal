package wrapper_layout

import (
	"github.com/Rafael24595/go-terminal/engine/app/draw"
	"github.com/Rafael24595/go-terminal/engine/app/pager"
	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/app/viewmodel"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable/line"
	"github.com/Rafael24595/go-terminal/engine/render/text"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

// TODO: Implement tokenize lines method to prevent line feed injection.
func TerminalApply(state *state.UIState, vm viewmodel.ViewModel, size terminal.Winsize) []text.Line {
	rows := int(size.Rows)
	cols := int(size.Cols)

	header, footer := vm.InitStaticLayers(size)

	headerLines := drawStaticLines(header.ToDrawable(), rows, cols)
	footerLines := drawStaticLines(footer.ToDrawable(), rows, cols)

	inputLines := make([]text.Line, 0)
	if input, ok := vm.InitInputLine(size); ok {
		inputLines = drawStaticLines(input, rows, cols)
	}

	helperLines := make([]text.Line, 0)
	if helper, ok := vm.InitHelper(size); ok {
		helperLines = drawStaticLines(helper, rows, cols)
	}

	static := len(headerLines) + len(footerLines) + len(inputLines) + len(helperLines)
	rest := int(size.Rows) - static
	if rest < 0 {
		return text.NewLines(
			text.LineFromString("Too low resolution"),
		)
	}

	remSize := terminal.NewWinsize(uint16(rest), size.Cols)
	lines := vm.InitDynamicLayers(remSize)

	dynamicSize := terminal.NewWinsize(uint16(rest), size.Cols)
	drawCtx := draw.NewDrawContext(state, dynamicSize)
	drawStt := drawDynamicLines(drawCtx, vm.Pager, lines.ToDrawable())

	state.Pager.Page = drawStt.Page
	state.Pager.HasMore = showPagination(drawStt)

	allLines := headerLines
	allLines = append(allLines, drawStt.Buffer...)
	allLines = append(allLines, footerLines...)
	allLines = append(allLines, inputLines...)
	allLines = append(allLines, helperLines...)

	return allLines
}

func drawStaticLines(drawable drawable.Drawable, rows, cols int) []text.Line {
	buffer := make([]text.Line, 0)

	content := true
	for content {
		lines, status := drawable.Draw()
		content = status

		if len(lines) == 0 {
			break
		}

		for _, lin := range lines {
			buffer = append(buffer,
				line.WrapLineWords(cols, lin)...,
			)

			if len(buffer) >= rows {
				break
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

	var lines []text.Line
	hasNext := true

	for hasNext {
		lines, hasNext = drawable.Draw()
		renderedSize := len(lines)
		if len(lines) == 0 {
			continue
		}

		state.Work.Reset()
		state.Work.Add(renderedSize)

		for _, lin := range lines {
			fixed := line.WrapLineWords(cols, lin)

			state.Work.Advance()
			state.Work.Add(len(fixed))

			for _, v := range fixed {
				state.Buffer[state.Cursor] = v

				state.Work.Advance()

				if f := text.HasFocus(v); f {
					state.Focus = f
				}

				state.Cursor += 1
				if state.Cursor < ctx.Size.Rows {
					continue
				}

				if shouldStop(ctx, pager, state) {
					return state
				}

				state = pager.Engine.Func(ctx, state)
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
	return stt.Page != 0 || (stt.Work.HasWorks() && stt.Work.Finished())
}
