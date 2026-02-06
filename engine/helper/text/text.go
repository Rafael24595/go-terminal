package text

import (
	"github.com/Rafael24595/go-terminal/engine/core/key"
	"github.com/Rafael24595/go-terminal/engine/helper/math"
)

var wrapperMap = map[rune]rune{
	'{': '}',
	'(': ')',
	'[': ']',
	'<': '>',
}

var runesRequiringTrailingSpace = []rune{
	',',
	'.',
	';',
}

var FullTextTransformer = NewTextTransformer(AppendSpaceAfterRune, WrappRunes)

type textTransform func(
	ky key.Key,
	start,
	end uint,
	buff []rune,
) (uint, uint, []rune, bool)

type TextTransformer struct {
	helpers []textTransform
}

func NewTextTransformer(helpers ...textTransform) TextTransformer {
	return TextTransformer{
		helpers: helpers,
	}
}

func (h TextTransformer) Apply(ky key.Key, start, end uint, buff []rune) (uint, uint, []rune) {
	for _, h := range h.helpers {
		if start, end, text, ok := h(ky, start, end, buff); ok {
			return start, end, text
		}
	}

	return start, end, []rune{ky.Rune}
}

func WrappRunes(ky key.Key, start, end uint, buff []rune) (uint, uint, []rune, bool) {
	text := []rune{ky.Rune}

	open := ky.Rune

	close, ok := wrapperMap[open]
	if !ok {
		return start, end, text, false
	}

	start = math.SubClampZero(start, 1)

	text = make([]rune, 0)
	text = append(text, open)
	text = append(text, buff[start:end]...)
	text = append(text, close)

	return start, end, text, true
}

func AppendSpaceAfterRune(ky key.Key, start, end uint, _ []rune) (uint, uint, []rune, bool) {
	text := []rune{ky.Rune}

	for _, r := range runesRequiringTrailingSpace {
		if ky.Rune != r {
			continue
		}

		text = append(text, ' ')

		return start, end, text, true
	}

	return start, end, text, false
}
