package headers

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	h := make(Headers)
	return h
}

const crlf = "\r\n"

func (h Headers) Parse(data []byte) (n int, done bool, err error) {

	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return 0, false, nil
	}

	if idx == 0 {
		return len(crlf), true, nil
	}

	line := data[:idx]
	parts := bytes.SplitN(line, []byte(":"), 2)

	if len(parts) < 2 {
		return 0, false, fmt.Errorf("malformed header (missing colon): %q", string(line))
	}

	key := strings.ToLower(string(parts[0]))

	if key != strings.TrimRight(key, " ") {
		return 0, false, fmt.Errorf("invalid header name: %s", key)
	}

	if !isValid(key) {
		return 0, false, fmt.Errorf("invalid use of characters %s", key)
	}

	key = strings.TrimSpace(key)
	val := string(bytes.TrimSpace(parts[1]))

	if oldVal, ok := h[key]; ok {
		joinedValue := oldVal + ", " + val
		val = joinedValue
	}

	h[key] = val

	return idx + len(crlf), false, nil

}

func (h Headers) Get(key string) (string, error) {
	val, ok := h[strings.ToLower(key)]
	if ok {
		return val, nil
	}
	return "", fmt.Errorf("Key does not exist")
}

func isValid(s string) bool {
	if len(s) < 1 {
		return false
	}

	for _, c := range s {
		switch {
		case c >= 'A' && c <= 'Z':
		case c >= 'a' && c <= 'z':
		case c >= '0' && c <= '9':
		case c == '!' || c == '#' || c == '$' || c == '%' || c == '&' || c == '\'' || c == '*' || c == '+' || c == '-' || c == '.' || c == '^' || c == '_' || c == '`' || c == '|' || c == '~':
		default:
			return false
		}
	}

	return true

}
