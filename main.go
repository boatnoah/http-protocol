package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const inputFilePath = "messages.txt"

func main() {
	file, err := os.Open("message.txt")

	if err != nil {
		log.Fatalf("could not open %s: %s\n", inputFilePath, err)
	}

	lines := getLinesChannel(file)

	for line := range lines {
		fmt.Printf("read: %s\n", line)
	}

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
