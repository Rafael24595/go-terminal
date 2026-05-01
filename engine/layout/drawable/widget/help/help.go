package help

import (
	"fmt"
	"strings"

	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/block"
	"github.com/Rafael24595/go-reacterm-core/engine/model/help"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

const Name = "help_drawable"

type HelpDrawable struct {
	loaded   bool
	meta     *help.HelpMeta
	drawable drawable.Drawable
}

func New(meta *help.HelpMeta) *HelpDrawable {
	return &HelpDrawable{
		loaded: false,
		meta:   meta,
	}
}

func DrawableFromMeta(meta *help.HelpMeta) drawable.Drawable {
	return New(meta).ToDrawable()
}

func (d *HelpDrawable) ToDrawable() drawable.Drawable {
	return drawable.Drawable{
		Name: Name,
		Code: d.drawable.Code,
		Tags: d.drawable.Tags,
		Init: d.init,
		Wipe: d.wipe,
		Draw: d.draw,
	}
}

func (d *HelpDrawable) init() {
	d.loaded = true

	d.drawable = makeDrawable(d.meta)

	d.drawable.Init()
}

func (d *HelpDrawable) wipe() {
	if d.drawable.Wipe == nil {
		return
	}
	d.drawable.Wipe()
}

func (d *HelpDrawable) draw(size winsize.Winsize) ([]text.Line, bool) {
	assert.True(d.loaded, drawable.MessageInitialized)

	return d.drawable.Draw(size)
}

func makeDrawable(meta *help.HelpMeta) drawable.Drawable {
	if len(meta.Fields) == 0 {
		return block.DrawableFromLines()
	}

	frags := make([]text.Fragment, len(meta.Fields))

	for i, field := range meta.Fields {
		code := strings.Join(field.Code, ", ")

		separator := ""
		if i < len(meta.Fields)-1 {
			separator = " | "
		}

		frag := fmt.Sprintf("[%s] %s%s", code, field.Detail, separator)
		frags = append(frags,
			*text.NewFragment(frag).
				AddAtom(style.AtmWrap),
		)
	}

	return block.DrawableFromLines(
		*text.EmptyLine(),
		*text.LineFromFragments(
			*text.NewFragment("--Help--"),
			*text.NewFragment("-").
				AddSpec(style.SpecFromKind(style.SpcKindFill)),
		),
		*text.LineFromFragments(frags...),
		*text.NewLine("-", style.SpecFromKind(style.SpcKindFill)),
	)
}
