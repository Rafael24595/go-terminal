package wrapper_screen

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/partial/pipeline"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/partial/pipeline/header"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

func NewBaseHeader() pipeline.Transformer {
	lines := text.ApplyLineSpec(
		style.SpecFromKind(style.SpcKindPaddingCenter),
		*text.LineFromFragments(
			*text.NewFragment("Lorem ipsum dolor sit amet").AddAtom(style.AtmUpper),
		),
		*text.LineFromFragments(
			*text.NewFragment("consectetur adipiscing").AddAtom(style.AtmUpper),
		),
		*text.LineFromFragments(
			*text.NewFragment("-Server 00-").AddAtom(style.AtmUpper),
		),
	)

	return header.HeaderTransformer(pipeline.Before, lines...)
}
