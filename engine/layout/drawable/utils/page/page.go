package page

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/draw"
	"github.com/Rafael24595/go-reacterm-core/engine/app/pager"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
	"github.com/Rafael24595/go-reacterm-core/engine/render/wrap"
)

func NewPageRenderer(strategy pager.PagerStrategy) draw.PageRenderer {
	return func(uiState *state.UIState, size winsize.Winsize, drawable drawable.Drawable) *draw.DrawState {
		ctx := draw.NewDrawContext(uiState, size)
		status := draw.NewDrawStatus(ctx)
		if size.Rows == 0 {
			return status
		}

		status.Work.Add(1)

		for status.Work.Unfinished() {
			status.Work.Advance()
			status.Work.Reset()

			lines, hasNext := drawable.Draw(size)
			if hasNext {
				status.Work.Add(1)
			}

			linesLen := uint(len(lines))
			if linesLen == 0 {
				continue
			}

			status.Work.Add(linesLen)

			for _, lne := range lines {
				fixed := wrap.WrapLine(ctx.Size.Cols, &lne)

				fixedLen := uint(len(fixed))
				if fixedLen == 0 {
					continue
				}

				status.Work.Advance()
				status.Work.Add(fixedLen)

				for _, fix := range fixed {
					status.SetAndNext(fix)
					status.Work.Advance()

					status.MarkFocus(
						text.HasAtom(style.AtmFocus, fix),
					)

					if !status.IsFull() {
						continue
					}

					if shouldStop(ctx, strategy, status) {
						return status
					}

					if status.Work.Unfinished() {
						status = strategy.Engine.Func(ctx, status)
					}
				}
			}
		}

		return status
	}
}

func shouldStop(
	ctx *draw.DrawContext,
	strategy pager.PagerStrategy,
	status *draw.DrawState,
) bool {
	args := pager.PredicateContext{
		Page:     status.Page,
		HasFocus: status.Focus,
	}
	return strategy.Predicate.Func(*ctx.State, args)
}
