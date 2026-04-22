package wrapper_reader

import (
	"bufio"
	"os"
	"strings"
	"time"
	"unicode"

	"github.com/Rafael24595/go-terminal/engine/model/ascii"
	"github.com/Rafael24595/go-terminal/engine/model/key"
)

const (
	ansiModShift     = "2"
	ansiModAlt       = "3"
	ansiModCtrl      = "5"
	ansiModCtrlShift = "6"
)

type KeyReader struct {
	reader *bufio.Reader
}

func New() *KeyReader {
	return &KeyReader{
		reader: bufio.NewReader(os.Stdin),
	}
}

func (r *KeyReader) ReadKey() (*key.Key, error) {
	rn, _, err := r.reader.ReadRune()
	if err != nil {
		return key.NewKeySpace(), err
	}

	if ky, ok := key.ControlKeyMap[rn]; ok {
		return ky, nil
	}

	if rn != ascii.ESC {
		sntz, _ := sanitizeRune(rn)
		return key.NewKeyRune(sntz), nil
	}

	return r.processEscapeSequence()
}

func (r *KeyReader) processEscapeSequence() (*key.Key, error) {
	r.waitForMoreInput()

	if r.reader.Buffered() == 0 {
		return key.NewKeyCode(key.ActionEsc), nil
	}

	rn, _, err := r.reader.ReadRune()
	if err != nil {
		return nil, err
	}

	if rn == '[' {
		return r.readCSISequence()
	}

	return r.resolveAltKey(rn), nil
}

func (r *KeyReader) readCSISequence() (*key.Key, error) {
	var params strings.Builder
	for {
		rn, _, err := r.reader.ReadRune()
		if err != nil {
			return nil, err
		}

		if isCSITerminator(rn) {
			return r.parseCSI(params.String(), rn), nil
		}

		params.WriteRune(rn)
	}
}

func (r *KeyReader) parseCSI(params string, rn rune) *key.Key {
	mod := r.parseModifier(params)

	if rn == ascii.TILDE {
		return r.resolveTildeCSI(params, rn, mod)
	}

	return r.resolveFinalCSI(rn, mod)
}

func (r *KeyReader) parseModifier(params string) key.ModMask {
	if !strings.Contains(params, ";") {
		return key.ModNone
	}

	parts := strings.Split(params, ";")
	if len(parts) < 2 {
		return key.ModNone
	}

	switch parts[1] {
	case ansiModShift:
		return key.ModShift
	case ansiModAlt:
		return key.ModAlt
	case ansiModCtrl:
		return key.ModCtrl
	case ansiModCtrlShift:
		return key.ModShift | key.ModCtrl
	}

	return key.ModNone
}

func (r *KeyReader) resolveTildeCSI(params string, rn rune, mod key.ModMask) *key.Key {
	if ky, ok := key.CsiTildeMap[params]; ok {
		return key.NewKeyCode(ky, mod)
	}

	sntz, _ := sanitizeRune(rn)
	return key.NewKeyRune(sntz)
}

func (r *KeyReader) resolveFinalCSI(rn rune, mod key.ModMask) *key.Key {
	if ky, ok := key.CsiFinalMap[rn]; ok {
		return key.NewKeyCode(ky, mod)
	}

	sntz, _ := sanitizeRune(rn)
	return key.NewKeyRune(sntz)
}

func (r *KeyReader) resolveAltKey(rn rune) *key.Key {
	exists, ok := key.AltKeyMap[rn]
	if ok {
		return exists
	}

	sntz, _ := sanitizeRune(rn)
	return key.NewKeyRune(sntz)
}

func (r *KeyReader) waitForMoreInput() {
	if r.reader.Buffered() == 0 {
		time.Sleep(1 * time.Millisecond)
	}
}

func isCSITerminator(rn rune) bool {
	return (rn >= 'A' && rn <= 'Z') || rn == ascii.TILDE
}

func isControlChar(rn rune) bool {
	return rn < ascii.SPACE || rn == ascii.DEL
}

func sanitizeRune(rn rune) (rune, bool) {
	if isControlChar(rn) || !unicode.IsPrint(rn) {
		return 0, false
	}
	return rn, true
}
