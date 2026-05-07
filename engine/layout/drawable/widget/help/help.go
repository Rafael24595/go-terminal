package help

import (
	"fmt"
	"strings"

	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/builder"
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

const Name = "help_drawable"

type HelpDrawable struct {
	loaded   bool
	fields   []key.Descriptor
	drawable drawable.Drawable
}

func New(fields []key.Descriptor) *HelpDrawable {
	return &HelpDrawable{
		loaded: false,
		fields: fields,
	}
}

func DrawableFromFields(fields []key.Descriptor) drawable.Drawable {
	return New(fields).ToDrawable()
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

	d.drawable = makeDrawable(d.fields)

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

func makeDrawable(fields []key.Descriptor) drawable.Drawable {
	fieldsLen := len(fields)
	if fieldsLen == 0 {
		return builder.DrainFromLines()
	}

	frags := make([]text.Fragment, fieldsLen)

	for i, field := range fields {
		code := strings.Join(field.Code, ", ")

		separator := ""
		if i < fieldsLen-1 {
			separator = " | "
		}

		frag := fmt.Sprintf("[%s] %s%s", code, field.Detail, separator)
		frags = append(frags,
			*text.NewFragment(frag).
				AddAtom(style.AtmWrap),
		)
	}

	return builder.DrainFromLines(
		*text.LineFromFragments(
			*text.NewFragment("--Help--"),
			*text.NewFragment("-").
				AddSpec(style.SpecFromKind(style.SpcKindFill)),
		),
		*text.LineFromFragments(frags...),
		*text.NewLine("-", style.SpecFromKind(style.SpcKindFill)),
	)
}
