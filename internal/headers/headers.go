package headers

import (
	"fmt"
	"strings"
)

func ValidateKey(key string) bool {
	//	Uppercase letters: A-Z
	//
	// Lowercase letters: a-z
	// Digits: 0-9
	// Special characters: !, #, $, %, &, ', *, +, -, ., ^, _, `, |, ~

	allowed := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!#$%&'*+-.^_`|~"
	for _, ch := range key {
		if !strings.ContainsRune(allowed, ch) {
			return false // invalid character found
		}
	}
	return true
}

const CRLF = "\r\n"

type Headers map[string]string

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	converted_str := string(data)
	idx := strings.Index(converted_str, CRLF)

	switch idx {
	case -1:
		return 0, false, nil
	case 0:
		return len(converted_str), true, nil
	default:
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
	if !ValidateKey(key_val[0]) {
		return 0, false, fmt.Errorf("field Name is not Valid")
	}
	key_val[1] = strings.TrimSpace(key_val[1])
	key_val[0] = strings.ToLower(key_val[0])
	if val, ok := h[key_val[0]]; ok && val != "" {
		h[key_val[0]] = val + "," + key_val[1]
	} else {
		h[key_val[0]] = key_val[1]
	}
	return idx + len(CRLF), false, nil
}
