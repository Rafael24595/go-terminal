package help

import (
	"fmt"
	"strings"

	"github.com/Rafael24595/go-terminal/engine/core/assert"
	"github.com/Rafael24595/go-terminal/engine/core/drawable"
	"github.com/Rafael24595/go-terminal/engine/core/drawable/line"
	"github.com/Rafael24595/go-terminal/engine/core/help"
	"github.com/Rafael24595/go-terminal/engine/core/style"
	"github.com/Rafael24595/go-terminal/engine/core/text"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

type HelpDrawable struct {
	initialized bool
	size        terminal.Winsize
	meta        *help.HelpMeta
	drawable    drawable.Drawable
}

func NewHelpDrawable(meta *help.HelpMeta) *HelpDrawable {
	return &HelpDrawable{
		initialized: false,
		size:        terminal.Winsize{},
		meta:        meta,
	}
}

func HelpDrawableFromMeta(meta *help.HelpMeta) drawable.Drawable {
	return NewHelpDrawable(meta).ToDrawable()
}

func (d *HelpDrawable) ToDrawable() drawable.Drawable {
	return drawable.Drawable{
		Init: d.init,
		Draw: d.draw,
	}
}

func (d *HelpDrawable) init(size terminal.Winsize) {
	d.initialized = true

	d.size = size

	d.drawable = makeDrawable(d.meta)

	d.drawable.Init(size)
}

func (d *HelpDrawable) draw() ([]text.Line, bool) {
	assert.True(d.initialized, "the drawable should be initialized before draw")

	return d.drawable.Draw()
}

func makeDrawable(meta *help.HelpMeta) drawable.Drawable {
	frags := make([]text.Fragment, len(meta.Fields))

	for i, field := range meta.Fields {
		code := strings.Join(field.Code, ", ")

		separator := ""
		if i < len(meta.Fields)-1 {
			separator = " | "
		}

		frag := fmt.Sprintf("[%s] %s%s", code, field.Detail, separator)
		frags = append(frags,
			text.NewFragment(frag).
				AddAtom(style.AtmWrap),
		)
	}

	return line.EagerDrawableFromLines(
		text.EmptyLine(),
		text.LineFromFragments(frags...),
	)
}
