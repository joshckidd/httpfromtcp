package response

import (
	"bytes"
	"fmt"
	"httpfromtcp/internal/headers"
	"net"
	"strconv"
)

type StatusCode int

const (
	S200 StatusCode = iota
	S400
	S500
)

type Writer struct {
	Conn       net.Conn
	Headers    headers.Headers
	StatusCode StatusCode
	Body       *bytes.Buffer
}

func (w *Writer) WriteStatusLine() error {
	var err error
	switch w.StatusCode {
	case S200:
		_, err = w.Conn.Write([]byte("HTTP/1.1 200 OK\r\n"))
	case S400:
		_, err = w.Conn.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
	case S500:
		_, err = w.Conn.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
	}
	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h["content-length"] = strconv.Itoa(contentLen)
	h["connection"] = "close"
	h["content-type"] = "text/plain"

	return h
}

func (w *Writer) WriteHeaders() error {
	w.Headers["content-length"] = strconv.Itoa(w.Body.Len())

	var err error
	for k, v := range w.Headers {
		_, err = w.Conn.Write([]byte(fmt.Sprintf("%s: %s\r\n", k, v)))
		if err != nil {
			break
		}
	}
	w.Conn.Write([]byte("\r\n"))
	return err
}

func (w *Writer) WriteBody() error {
	_, err := w.Body.WriteTo(w.Conn)
	return err
}

func (w *Writer) Write() error {
	err := w.WriteStatusLine()
	if err != nil {
		return err
	}
	err = w.WriteHeaders()
	if err != nil {
		return err
	}
	err = w.WriteBody()
	return err
}
