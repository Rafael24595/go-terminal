package wrapper_terminal

import (
	"bufio"
	"os"
	"strings"
	"unicode"

	"github.com/Rafael24595/go-terminal/engine/core/key"
)

var controlKeyMap = map[rune]*key.Key{
	key.CTRL_C:   key.NewKeyCode(key.KeyCtrlC),
	key.CTRL_W:   key.NewKeyCode(key.KeyDeleteWord),
	key.TAB:      key.NewKeyCode(key.KeyTab),
	key.ENTER_LF: key.NewKeyCode(key.KeyEnter),
	key.ENTER_CR: key.NewKeyCode(key.KeyEnter),
	key.DEL:      key.NewKeyCode(key.KeyBackspace),
	key.BS:       key.NewKeyCode(key.KeyBackspace),
}

var altKeyMap = map[rune]*key.Key{
	'd': key.NewKeyCode(key.KeyDeleteWordForward, key.ModAlt),
}

var csiFinalMap = map[rune]key.KeyCode{
	'A': key.KeyArrowUp,
	'B': key.KeyArrowDown,
	'C': key.KeyArrowRight,
	'D': key.KeyArrowLeft,
	'H': key.KeyHome,
	'F': key.KeyEnd,
}

var csiTildeMap = map[string]key.KeyCode{
	"3": key.KeyDelete,
	"1": key.KeyHome,
	"7": key.KeyHome,
	"4": key.KeyEnd,
	"8": key.KeyEnd,
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

	exists, ok := controlKeyMap[char]
	if ok {
		return exists, nil
	}

	if char == key.ESC {
		return r.readEscapeSequence()
	}

	sntz, _ := sanitizeRune(char)
	return key.NewKeyRune(sntz), nil
}

func (r *inputReader) readEscapeSequence() (*key.Key, error) {
	char, _, err := r.reader.ReadRune()
	if err != nil {
		return nil, err
	}

	if char != '[' {
		return r.decodeAltKey(char), nil
	}

	return r.readCSISequence()
}

func (r *inputReader) decodeAltKey(char rune) *key.Key {
	exists, ok := altKeyMap[char]
	if ok {
		return exists
	}

	sntz, _ := sanitizeRune(char)
	return key.NewKeyRune(sntz)
}

func (r *inputReader) readCSISequence() (*key.Key, error) {
	params := ""
	for {
		ch, _, err := r.reader.ReadRune()
		if err != nil {
			return nil, err
		}

		if (ch >= 'A' && ch <= 'Z') || ch == '~' {
			return decodeCSI(params, ch), nil
		}
		params += string(ch)
	}
}

func decodeCSI(params string, char rune) *key.Key {
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

	if char == key.TILDE {
		return decodeTildeCSI(params, char, mod)
	}

	return decodeFinalCSI(char, mod)
}

func decodeTildeCSI(params string, char rune, mod key.ModMask) *key.Key {
	exists, ok := csiTildeMap[params]
	if ok {
		return key.NewKeyCode(exists, mod)
	}

	sntz, _ := sanitizeRune(char)
	return key.NewKeyRune(sntz)
}

func decodeFinalCSI(char rune, mod key.ModMask) *key.Key {
	exists, ok := csiFinalMap[char]
	if ok {
		return key.NewKeyCode(exists, mod)
	}

	sntz, _ := sanitizeRune(char)
	return key.NewKeyRune(sntz)
}

func sanitizeRune(r rune) (rune, bool) {
	if r < 0x20 || r == 0x7f {
		return 0, false
	}

	if !unicode.IsPrint(r) {
		return 0, false
	}

	return r, true
}
