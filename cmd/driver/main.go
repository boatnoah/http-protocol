package main

import (
	"github.com/httpfromtcp/internal/request"
	"strings"
)

func main() {

	_, err := request.RequestFromReader(strings.NewReader("GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))

	if err != nil {
		return
	}

}
