package wrapper_terminal

import (
	"bufio"
	"os"
	"strings"

	"github.com/Rafael24595/go-terminal/engine/core/key"
)

var csiFinalMap = map[rune]key.KeyCode{
	'A': key.KeyArrowUp,
	'B': key.KeyArrowDown,
	'C': key.KeyArrowRight,
	'D': key.KeyArrowLeft,
}

type inputReader struct {
	reader *bufio.Reader
}

func newInputReader() *inputReader {
	stdin := os.Stdin
	return &inputReader{
		reader: bufio.NewReader(stdin),
	}
}

func (r *inputReader) readRune() (*key.Key, error) {
	char, _, err := r.reader.ReadRune()
	if err != nil {
		return key.NewKeySpace(), err
	}

	switch char {
	case key.CTRL_C:
		return key.NewKeyCode(key.KeyCtrlC), nil
	case key.TAB:
		return key.NewKeyCode(key.KeyTab), nil
	case key.ENTER_LF, key.ENTER_CR:
	return key.NewKeyCode(key.KeyEnter), nil
	case key.DEL, key.BS:
		return key.NewKeyCode(key.KeyBackspace), nil
	case key.ESC:
		return r.readEscRune()
	default:
		return key.NewKeyRune(char), nil
	}
}

func (r *inputReader) readEscRune() (*key.Key, error) {
	next, _, err := r.reader.ReadRune()
	if err != nil {
		return nil, err
	}

	if next != '[' {
		return key.NewKeyRune(next), nil
	}

	params := ""
	for {
		ch, _, err := r.reader.ReadRune()
		if err != nil {
			return nil, err
		}

		if ch >= 'A' && ch <= 'Z' {
			return decodeCSI(params, ch), nil
		}
		params += string(ch)
	}
}

func decodeCSI(params string, final rune) *key.Key {
	mod := key.ModNone

	if strings.Contains(params, ";") {
		parts := strings.Split(params, ";")
		if len(parts) == 2 {
			switch parts[1] {
			case "2":
				mod = key.ModShift
			case "3":
				mod = key.ModAlt
			case "5":
				mod = key.ModCtrl
			case "6":
				mod = key.ModShift | key.ModCtrl
			}
		}
	}

	exists, ok := csiFinalMap[final]
	if !ok {
		return key.NewKeyRune(final)
	}

	return key.NewKeyCode(exists, mod)
}
