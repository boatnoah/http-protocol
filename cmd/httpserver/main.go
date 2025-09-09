package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/httpfromtcp/internal/headers"
	"github.com/httpfromtcp/internal/request"
	"github.com/httpfromtcp/internal/response"
	"github.com/httpfromtcp/internal/server"
)

const port = 42069

func routerHandler(w io.Writer, req *request.Request) *server.HandlerError {
	rw := response.NewWriter(w)

	// Common HTML bodies
	const html400 = `
<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`
	const html500 = `
<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`
	const html200 = `
<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`

	h := headers.NewHeaders()
	h.Set("content-type", "text/html")

	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		_ = rw.WriteStatusLine(response.BadRequest)
		_ = rw.WriteHeaders(h)
		_, _ = rw.WriteBody([]byte(html400))
		return nil

	case "/myproblem":
		_ = rw.WriteStatusLine(response.InternalServerError)
		_ = rw.WriteHeaders(h)
		_, _ = rw.WriteBody([]byte(html500))
		return nil

	default:
		_ = rw.WriteStatusLine(response.Ok)
		_ = rw.WriteHeaders(h)
		_, _ = rw.WriteBody([]byte(html200))
		return nil
	}
}

func main() {
	server, err := server.Serve(routerHandler, port)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
