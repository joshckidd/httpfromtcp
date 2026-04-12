package request

import (
	"errors"
	"io"
	"regexp"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	req, err := io.ReadAll(reader)
	if err != nil {
		return &Request{}, err
	}
	return parseRequestLine(strings.Split(string(req), "\r\n")[0])
}

func parseRequestLine(line string) (*Request, error) {
	parts := strings.Split(line, " ")

	if len(parts) > 3 {
		return &Request{}, errors.New("Invalid number of parts.")
	}
	r, _ := regexp.Compile("[^A-Z]+")
	if r.MatchString(parts[0]) {
		return &Request{}, errors.New("Method must be only capital alphabetic characters.")
	}
	var version string
	if strings.Contains(parts[2], "/") {
		version = strings.Split(parts[2], "/")[1]
		if version != "1.1" {
			return &Request{}, errors.New("Invalid http version.")
		}
	} else {
		return &Request{}, errors.New("Invalid http version.")
	}

	return &Request{
		RequestLine: RequestLine{
			Method:        parts[0],
			RequestTarget: parts[1],
			HttpVersion:   version,
		},
	}, nil
}
