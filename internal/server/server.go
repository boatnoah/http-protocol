package server

import (
	"fmt"
	"io"
	"net"
	"sync/atomic"

	"github.com/httpfromtcp/internal/headers"
	"github.com/httpfromtcp/internal/request"
	"github.com/httpfromtcp/internal/response"
)

type Server struct {
	listener net.Listener
	closed   atomic.Bool
	handler  Handler
}

type HandlerError struct {
	statusCode response.StatusCode
	message    string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

func Serve(handler Handler, port int) (*Server, error) {
	addr := fmt.Sprintf(":%d", port)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("Error creating listener %s", err)
	}
	s := &Server{listener: l, handler: handler}
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

	req, err := request.RequestFromReader(conn)
	if err != nil {
		(&HandlerError{statusCode: response.BadRequest}).Write(conn)
		return
	}

	if s.handler == nil {
		(&HandlerError{statusCode: response.InternalServerError}).Write(conn)
		return
	}

	if he := s.handler(conn, req); he != nil {
		he.Write(conn)
	}
}

func (he HandlerError) Write(w io.Writer) {
	rw := response.NewWriter(w)
	_ = rw.WriteStatusLine(he.statusCode)

	h := headers.NewHeaders()
	h.Set("content-type", "text/html")
	_ = rw.WriteHeaders(h)

	var body []byte
	switch he.statusCode {
	case response.BadRequest:
		body = []byte(`
<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`)
	case response.InternalServerError:
		body = []byte(`
<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`)
	default:
		body = []byte(`
<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`)
	}
	_, _ = rw.WriteBody(body)
}
