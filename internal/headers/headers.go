package headers

import (
	"fmt"
	"strings"
)

func TrimPrefixWhiteSpace(str string) string {
	for idx, ch := range str {
		if ch == ' ' {
			continue
		} else {
			str = str[idx:]
			break
		}
	}
	return str
}

func TrimSuffixWhiteSpace(str string) string {
	for idx, ch := range str {
		if ch == ' ' {
			str = str[:idx]
			break
		}
	}
	return str
}

const CRLF = "\r\n"

type Headers map[string]string

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	converted_str := string(data)
	idx := strings.Index(converted_str, CRLF)
	if idx == -1 {
		return 0, false, nil
	} else if idx == 0 {
		return len(converted_str), true, nil
	} else {
		converted_str = converted_str[:idx]
	}

	converted_str = strings.TrimSpace(converted_str)

	key_val := strings.SplitN(converted_str, ":", 2)
	if len(key_val) < 2 {
		return 0, false, fmt.Errorf("not enough arguments")
	}
	if strings.Contains(key_val[0], " ") {
		return 0, false, fmt.Errorf("string contains space between field val and colon")
	}

	key_val[1] = strings.TrimSpace(key_val[1])

	h[key_val[0]] = key_val[1]
	return idx + len(CRLF), false, nil
}
