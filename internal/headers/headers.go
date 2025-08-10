package headers

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"

	"github.com/MrBhop/httpfromtcp/internal/request"
)

type Headers map[string]string

func NewHeaders() Headers {
	return Headers{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	n = strings.Index(string(data), request.CrLf)
	switch n {
	case -1:
		return 0, false, nil
	case 0:
		return 2, true, nil
	}

	line := data[:n]
	parts := bytes.SplitN(line, []byte(":"), 2)

	key := strings.ToLower(string(parts[0]))
	if key != strings.TrimRight(key, " ") {
		return 0, false, fmt.Errorf("Whitespace after the key is not allowed.")
	}

	key = strings.TrimSpace(key)
	if err := validateHeaderKey(key); err != nil {
		return 0, false, err
	}

	fieldValue := strings.TrimSpace(string(parts[1]))
	if len(strings.Split(fieldValue, " ")) > 1 {
		return 0, false, fmt.Errorf("Whitespace in the feld value is not allowed.")
	}

	h[key] = fieldValue
	return n + 2, false, nil
}

func validateHeaderKey(key string) error {
	if len(key) < 1 {
		return fmt.Errorf("Invalid key length")
	}

	specialChars := map[rune]struct{}{
		'!': {},
		'#': {},
		'$': {},
		'%': {},
		'&': {},
		'\'': {},
		'*': {},
		'+': {},
		'-': {},
		'.': {},
		'^': {},
		'_': {},
		'`': {},
		'|': {},
		'~': {},
	}
	for _, r := range key {
		if unicode.IsLetter(r) {
			continue
		}

		if unicode.IsDigit(r) {
			continue
		}

		if _, exists := specialChars[r]; exists {
			continue
		}

		return fmt.Errorf("Invalid character")
	}

	return nil
}
