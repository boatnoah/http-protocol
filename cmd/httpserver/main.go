package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/httpfromtcp/internal/request"
	"github.com/httpfromtcp/internal/response"
	"github.com/httpfromtcp/internal/server"
)

const port = 42069

func routerHandler(w io.Writer, req *request.Request) *server.HandlerError {
	path := req.RequestLine.RequestTarget

	switch path {
	case "/yourproblem":
		return server.NewHandlerError(response.BadRequest, "Your problem is not my problem\n")
	case "/myproblem":
		return server.NewHandlerError(response.InternalServerError, "Woopsie, my bad\n")
	default:
		_, _ = io.WriteString(w, "All good, frfr\n")
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
