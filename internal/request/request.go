package request

import (
	"bytes"
	"errors"
	"httpfromtcp/internal/headers"
	"io"
	"regexp"
	"strings"
)

type State int

const (
	initialized State = iota
	done
	parsingHeaders
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	State       State
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	req := Request{
		State:   initialized,
		Headers: headers.NewHeaders(),
	}
	buf := make([]byte, 8)
	bytesRead := 0

	for req.State != done {
		n, err := reader.Read(buf[bytesRead:])
		if n == 0 && err == io.EOF {
			req.State = done
			break
		}
		if err != nil {
			return &req, err
		}
		bytesRead += n

		if bytesRead == len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		m, err := req.parse(buf[:bytesRead])
		if err != nil {
			return &req, err
		}

		if m != 0 {
			newBuf := make([]byte, len(buf))
			copy(newBuf, buf[m:])
			buf = newBuf
			bytesRead -= m
		}
	}

	return &req, nil
}

func parseRequestLine(req []byte) (*Request, int, error) {
	lines := bytes.Split(req, []byte("\r\n"))
	if len(lines) == 1 {
		return &Request{}, 0, nil
	}

	lineBytes := lines[0]
	n := len(lineBytes) + 2
	line := string(lineBytes)

	parts := strings.Split(line, " ")

	if len(parts) < 3 {
		return &Request{}, n, errors.New("Invalid number of parts.")
	}

	if len(parts) > 3 {
		return &Request{}, n, errors.New("Invalid number of parts.")
	}
	r, _ := regexp.Compile("[^A-Z]+")
	if r.MatchString(parts[0]) {
		return &Request{}, n, errors.New("Method must be only capital alphabetic characters.")
	}
	var version string
	if strings.Contains(parts[2], "/") {
		version = strings.Split(parts[2], "/")[1]
		if version != "1.1" {
			return &Request{}, n, errors.New("Invalid http version.")
		}
	} else {
		return &Request{}, n, errors.New("Invalid http version.")
	}

	return &Request{
		RequestLine: RequestLine{
			Method:        parts[0],
			RequestTarget: parts[1],
			HttpVersion:   version,
		},
	}, n, nil
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.State {
	case initialized:
		req, n, err := parseRequestLine(data)

		if n != 0 {
			r.State = parsingHeaders
			r.RequestLine = req.RequestLine
		}
		return n, err
	case parsingHeaders:
		n, d, err := r.Headers.Parse(data)
		if d {
			r.State = done
		}

		return n, err
	}
	return 0, nil
}
