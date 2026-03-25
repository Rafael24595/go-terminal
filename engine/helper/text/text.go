package text

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
var VoidTextTransformer = NewTextTransformer()

type textTransform func(
	text []rune,
	start,
	end uint,
	buff []rune,
) ([]rune, bool)

type TextTransformer struct {
	helpers []textTransform
}

func NewTextTransformer(helpers ...textTransform) TextTransformer {
	return TextTransformer{
		helpers: helpers,
	}
}

func (h TextTransformer) Apply(text []rune, start, end uint, buff []rune) []rune {
	for _, h := range h.helpers {
		if text, ok := h(text, start, end, buff); ok {
			return text
		}
	}

	return text
}

func WrappRunes(text []rune, start, end uint, buff []rune) ([]rune, bool) {
	size := len(text)
	if size < 1 || size > 1 {
		return text, false
	}

	focus := text[0]

	close, ok := wrapperMap[focus]
	if !ok {
		return text, false
	}

	text = make([]rune, 0)
	text = append(text, focus)
	text = append(text, buff[start:end]...)
	text = append(text, close)

	return text, true
}

func AppendSpaceAfterRune(text []rune, start, end uint, _ []rune) ([]rune, bool) {
	size := len(text)
	if size < 1 || size > 1 {
		return text, false
	}

	focus := text[0]
	for _, r := range runesRequiringTrailingSpace {
		if focus != r {
			continue
		}

		text = append(text, ' ')

		return text, true
	}

	return text, false
}
