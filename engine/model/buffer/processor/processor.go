package processor

import (
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
)

type Processor func([]rune) ([]rune, []rune)

func Identity(buffer []rune) ([]rune, []rune) {
	return buffer, buffer
}

func Number(buffer []rune) ([]rune, []rune) {
	var fixedBuffer []rune
	hasSeparator := false

	for i, v := range buffer {
		if v >= '0' && v <= '9' {
			fixedBuffer = append(fixedBuffer, v)
			continue
		}

		if (v == ',' || v == '.') && !hasSeparator {
			fixedBuffer = append(fixedBuffer, v)
			hasSeparator = true
			continue
		}

		if v == '-' && len(fixedBuffer) == 0 {
			fixedBuffer = append(fixedBuffer, v)
		}

		if fixedBuffer == nil {
			fixedBuffer = make([]rune, 0, len(buffer))
			fixedBuffer = append(fixedBuffer, buffer[:i]...)
		}
	}

	if fixedBuffer == nil {
		return buffer, buffer
	}

	return fixedBuffer, fixedBuffer
}

func Hidden(buffer []rune) ([]rune, []rune) {
	fixedBuffer := make([]rune, len(buffer))

	for i, r := range buffer {
		switch r {
		case '\n', '\t', ' ':
			fixedBuffer[i] = r
		default:
			fixedBuffer[i] = '*'
		}
	}

	return buffer, fixedBuffer
}

func Inline(buffer []rune) ([]rune, []rune) {
	fixedBuffer := make([]rune, len(buffer))

	for i, r := range buffer {
		switch r {
		case '\n':
			fixedBuffer[i] = ' '
		default:
			fixedBuffer[i] = buffer[i]
		}
	}

	return fixedBuffer, fixedBuffer
}

func Limit(limit winsize.Cols, processor Processor) Processor {
	if limit == 0 {
		return processor
	}

	return func(buffer []rune) ([]rune, []rune) {
		buffer, facade := processor(buffer)
		return trimBuffer(limit, facade, buffer)
	}
}

func trimBuffer(limit winsize.Cols, facade, buffer []rune) ([]rune, []rune) {
	if limit == 0 {
		return buffer, facade
	}

	if len(buffer) > int(limit) {
		buffer = buffer[:limit]
	}

	if len(facade) > int(limit) {
		facade = facade[:limit]
	}

	return buffer, facade
}
