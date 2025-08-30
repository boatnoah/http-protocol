package main

import (
	"fmt"
	"github.com/httpfromtcp/internal/request"
	"log"
	"net"
)

func main() {

	ln, err := net.Listen("tcp", ":42069")

	if err != nil {
		log.Fatalf("failed to connect to port %s, %s", ":42069", err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatalf("failed to accept, %s", err)
			break
		}
		r, err := request.RequestFromReader(conn)

		if err != nil {
			log.Fatalf("failed to request from reader %s", err)
		}

		fmt.Println("Request line:")
		fmt.Println("- Method:", r.RequestLine.Method)
		fmt.Println("- Target:", r.RequestLine.RequestTarget)
		fmt.Println("- Version:", r.RequestLine.HttpVersion)

	}

	fmt.Println("==Connection has been closed==")

}
