package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
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

		fmt.Println("==Connection has been accepted==")

		lines := getLinesChannel(conn)

		for line := range lines {
			fmt.Printf("read: %s\n", line)
		}
	}

	fmt.Println("==Connection has been closed==")

}
func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)

	go func() {

		defer close(ch)
		defer f.Close()
		currLine := ""
		for {
			buffer := make([]byte, 8, 8)
			n, err := f.Read(buffer)
			if err != nil {
				return
			}
			parts := strings.Split(string(buffer[:n]), "\n")

			for i, part := range parts {
				if i == len(parts)-1 {
					break
				}
				line := fmt.Sprintf("%s%s", currLine, part)
				ch <- line
				currLine = ""
			}

			currLine += parts[len(parts)-1]
		}

	}()

	return ch
}
