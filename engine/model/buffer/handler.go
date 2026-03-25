package buffer

type InputType uint8

const (
	String InputType = iota
	Number
	Hidden
)

type RuneHandler func([]rune) ([]rune, []rune)

func voidRuneHandler(buffer []rune) ([]rune, []rune) {
	return buffer, buffer
}

func NewRuneHandler(input InputType) RuneHandler {
	return NewLimitedRuneHandler(0, input)
}

func NewLimitedRuneHandler(limit uint64, input InputType) RuneHandler {
	if limit == 0 && input == String {
		return voidRuneHandler
	}

	return func(buff []rune) ([]rune, []rune) {
		var facade []rune
		var buffer []rune

		switch input {
		case Number:
			buffer, facade = fixNumberBuffer(buff)
		case Hidden:
			buffer, facade = fixHiddenBuffer(buff)
		default:
			buffer, facade = fixStringBuffer(buff)
		}

		return trimBuffer(limit, facade, buffer)
	}
}

func fixStringBuffer(buff []rune) ([]rune, []rune) {
	return buff, buff
}

func fixNumberBuffer(buff []rune) ([]rune, []rune) {
	var fixedBuffer []rune
	hasSeparator := false

	for i, v := range buff {
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
			fixedBuffer = make([]rune, 0, len(buff))
			fixedBuffer = append(fixedBuffer, buff[:i]...)
		}
	}

	if fixedBuffer == nil {
		return buff, buff
	}

	return fixedBuffer, fixedBuffer
}

func fixHiddenBuffer(buff []rune) ([]rune, []rune) {
	fixedBuffer := make([]rune, len(buff))
	for i, r := range buff {
		switch r {
		case '\n', '\t', ' ':
			fixedBuffer[i] = r
		default:
			fixedBuffer[i] = '*'
		}
	}
	return buff, fixedBuffer
}

func trimBuffer(limit uint64, facade []rune, buffer []rune) ([]rune, []rune) {
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
