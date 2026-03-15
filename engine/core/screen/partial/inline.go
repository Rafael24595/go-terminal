package partial

import (
	"github.com/Rafael24595/go-terminal/engine/commons/structure/set"
	"github.com/Rafael24595/go-terminal/engine/core/drawable"
	"github.com/Rafael24595/go-terminal/engine/core/drawable/action"
	"github.com/Rafael24595/go-terminal/engine/core/inline"
	"github.com/Rafael24595/go-terminal/engine/core/screen"
	"github.com/Rafael24595/go-terminal/engine/core/screen/primitive"

	drawable_inline "github.com/Rafael24595/go-terminal/engine/core/drawable/inline"
)

type Inline struct {
	separator string
	screen    screen.Screen
	actions   []action.Action
}

func NewInline(screen screen.Screen) *Inline {
	return &Inline{
		separator: inline.DefaultInlineSeparator,
		screen:    screen,
		actions:   make([]action.Action, 0),
	}
}

func (c *Inline) Separator(separator string) *Inline {
	c.separator = separator
	return c
}

func (c *Inline) PushAction(focus action.Focus, filter inline.FilterMeta) *Inline {
	action := action.NewAction(
		action.ActionMapGroup,
		focus,
		c.groupDrawables(filter),
	)

	c.actions = append(c.actions, action)
	return c
}

func (c *Inline) ToScreen() screen.Screen {
	return primitive.NewMapScreen(c.screen).
		PushAction(c.actions...).
		ToScreen()
}

func (c *Inline) groupDrawables(filter inline.FilterMeta) action.ActionFunc {
	return func(drawables ...drawable.Drawable) []drawable.Drawable {
		rest := make([]drawable.Drawable, 0, len(drawables))
		filt := make([]drawable.Drawable, 0, len(drawables))

		values := set.SetFrom(filter.Values...)

		for _, d := range drawables {
			if filter.Target == inline.TargetCode && values.Has(d.Code) {
				filt = append(filt, d)
				continue
			}

			if filter.Target == inline.TargetTags && d.Tags.Any(values) {
				filt = append(filt, d)
				continue
			}

			rest = append(rest, d)
		}

		if len(filt) == 0 {
			return drawables
		}

		return append(rest,
			drawable_inline.NewInlineDrawable(filt...).
				Separator(c.separator).
				ToDrawable())
	}
}
