package focus

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/draw"
	"github.com/Rafael24595/go-reacterm-core/engine/app/pager"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/helper/math"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

// TODO: Add flag to manage non focus drawable?
func FocusInitTransformer(engine pager.Engine) pipeline.InitTransformer {
	return func(size winsize.Winsize, drawable drawable.Drawable) ([]text.Line, bool) {
		ctx := draw.NewDrawContext(state.NewUIState(), size)
		state := draw.NewDrawStatus(ctx)

		for {
			lines, status := drawable.Draw(size)

			for len(lines) > 0 {
				remaining := math.SubClampZero(
					size.Rows,
					winsize.Rows(state.Cursor),
				)

				limit := winsize.Rows(len(lines))
				if len(lines) > int(remaining) {
					limit = remaining
				}

				chunk := lines[:limit]
				lines = lines[limit:]

				for _, l := range chunk {
					state.SetAndNext(l)
				}

				//TODO: Remove focus atom to prevent conflicts?
				state.MarkFocus(
					text.HasAtom(style.AtmFocus, chunk...),
				)

				if state.Cursor == uint16(size.Rows) {
					if state.Focus {
						return state.Buffer, false
					}

					state = engine.Func(ctx, state)
				}
			}

			if !status {
				return state.Buffer, false
			}
		}
	}
}
