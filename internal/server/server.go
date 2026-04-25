package server

import (
	"bytes"
	"fmt"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"net"
	"sync/atomic"
)

type Server struct {
	Listener net.Listener
	Open     atomic.Bool
	Handler  Handler
}

type Handler func(w *response.Writer, req *request.Request)

func Serve(port int, h Handler) (*Server, error) {
	a := fmt.Sprintf(":%d", port)
	l, err := net.Listen("tcp", a)
	if err != nil {
		return &Server{}, err
	}

	s := Server{
		Listener: l,
		Handler:  h,
	}
	s.Open.Store(true)
	go s.listen()
	return &s, nil
}

func (s *Server) Close() error {
	s.Open.Store(false)
	return s.Listener.Close()
}

func (s *Server) listen() {
	for s.Open.Load() {
		conn, err := s.Listener.Accept()
		if err != nil {
			break
		}

		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	buf := bytes.NewBuffer([]byte(""))

	req, err := request.RequestFromReader(conn)
	if err != nil {
		buf.Write([]byte(err.Error()))
		rw := &response.Writer{
			Conn:       conn,
			Headers:    response.GetDefaultHeaders(buf.Len()),
			Body:       buf,
			StatusCode: response.S500,
		}
		rw.Write()
		return
	}

	s.Handler(&response.Writer{
		Conn:    conn,
		Headers: response.GetDefaultHeaders(0),
		Body:    buf,
	}, req)

	buf.WriteTo(conn)
	conn.Close()
}
