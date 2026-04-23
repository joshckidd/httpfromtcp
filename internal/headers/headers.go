package headers

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	lines := bytes.Split(data, []byte("\r\n"))
	if len(lines) == 1 {
		return 0, false, nil
	}

	lineBytes := lines[0]
	if len(lineBytes) == 0 {
		return 2, true, nil
	}

	l := len(lineBytes) + 2
	line := string(lineBytes)

	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return 0, false, errors.New("Header line does not contain :.")
	}

	if strings.TrimSpace(parts[0]) != parts[0] {
		return 0, false, errors.New("There must be no whitespace in the field name.")
	}

	r, _ := regexp.Compile("^[a-zA-Z0-9!#$%&'*+\\-.^_`|~]+$")

	if !r.MatchString(parts[0]) {
		return 0, false, errors.New("Invalid characters in header field.")
	}

	val, ok := h[strings.ToLower(parts[0])]
	if ok {
		val = fmt.Sprintf("%s, %s", val, strings.TrimSpace(parts[1]))
	} else {
		val = strings.TrimSpace(parts[1])
	}

	h[strings.ToLower(parts[0])] = val

	return l, false, nil
}

func (h Headers) Get(key string) (string, error) {
	var err error
	err = nil
	v, ok := h[strings.ToLower(key)]
	if !ok {
		err = errors.New("Invalid header key.")
	}

	return v, err
}
