package main

import (
	"github.com/httpfromtcp/internal/headers"
)

func main() {

	h := headers.NewHeaders()

	data1 := []byte("       Host : localhost:42069       \r\n\r\n")

	_, _, _ = h.Parse(data1)

}
