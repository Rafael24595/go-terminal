package wrapper_terminal

import (
	"bufio"
	"fmt"
	"os"
)

type inputReader struct {
	reader *bufio.Reader
}

func newInputReader() *inputReader {
	stdin := os.Stdin
	return &inputReader{
		reader: bufio.NewReader(stdin),
	}
}

func (r *inputReader) readRune() (string, error) {
	ch, _, err := r.reader.ReadRune()
	if err != nil {
		return "", err
	}

	if ch != 0x1b {
		return string(ch), nil
	}

	next, _, err := r.reader.ReadRune()
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	if next != '[' {
		return string([]rune{ch, next}), nil
	}

	final, _, err := r.reader.ReadRune()
	if err != nil {
		return "", err
	}

	switch final {
	case 'A':
		return ARROW_UP, nil
	case 'B':
		return ARROW_DOWN, nil
	case 'C':
		return ARROR_RIGHT, nil
	case 'D':
		return ARROW_LEFT, nil
	default:
		return UNKNOWN, nil
	}
}
