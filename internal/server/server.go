package server

import (
	"fmt"
	"net"
	"sync/atomic"

	"github.com/httpfromtcp/internal/response"
)

type Server struct {
	listener net.Listener
	closed   atomic.Bool
}

func Serve(port int) (*Server, error) {
	addr := fmt.Sprintf(":%d", port)
	l, err := net.Listen("tcp", addr)

	if err != nil {
		return nil, fmt.Errorf("Error creating listener %s", err)
	}

	s := &Server{
		listener: l,
	}

	go s.listen()

	return s, nil
}

func (s *Server) Close() error {
	if s.closed.Swap(true) {
		return nil
	}
	return s.listener.Close()
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			continue
		}
		go s.handle(conn)

	}

}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	h := response.GetDefaultHeaders(0)
	response.WriteStatusLine(conn, response.Ok)
	response.WriteHeaders(conn, h)
}
