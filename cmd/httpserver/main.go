package main

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
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

	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/") {
		path := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")
		up := "https://httpbin.org/" + path

		resp, err := http.Get(up)
		if err != nil {
			_ = rw.WriteStatusLine(response.InternalServerError)
			_ = rw.WriteHeaders(headers.NewHeaders())
			_, _ = rw.WriteBody([]byte("upstream error\n"))
			return nil
		}
		defer resp.Body.Close()

		_ = rw.WriteStatusLine(response.StatusCode(resp.StatusCode))

		h := headers.NewHeaders()
		if ct := resp.Header.Get("Content-Type"); ct != "" {
			h.Set("content-type", ct)
		}
		h.Set("transfer-encoding", "chunked")
		h.Set("trailer", "X-Content-Sha256, X-Content-Length")
		_ = rw.WriteHeaders(h)

		hasher := sha256.New()
		var total int64

		buf := make([]byte, 1024)
		for {
			n, rerr := resp.Body.Read(buf)
			if n > 0 {
				chunk := buf[:n]
				_, err = hasher.Write(chunk)
				if err != nil {
					return nil
				}
				total += int64(n)

				if _, err := rw.WriteChunkedBody(buf[:n]); err != nil {
					log.Printf("client closed during chunk write: %v", err)
					return nil
				}
			}
			if rerr == io.EOF {
				break
			}
			if rerr != nil {
				break
			}
		}

		sum := hasher.Sum(nil)
		tr := headers.NewHeaders()
		tr.Set("X-Content-SHA256", hex.EncodeToString(sum))
		tr.Set("X-Content-Length", strconv.FormatInt(total, 10))

		if err := rw.WriteTrailers(tr); err != nil {
			return nil
		}
		return nil
	}

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
	case "/video":
		data, err := os.ReadFile("assets/vim.mp4")
		if err != nil {
			_ = rw.WriteStatusLine(response.InternalServerError)
			_ = rw.WriteHeaders(h)
			_, _ = rw.WriteBody([]byte(html500))
			return nil
		}
		v := headers.NewHeaders()
		v.Set("content-type", "video/mp4")
		_ = rw.WriteStatusLine(response.Ok)
		_ = rw.WriteHeaders(v)
		_, _ = rw.WriteBody([]byte(data))
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
