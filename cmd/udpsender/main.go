package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {

	udp, err := net.ResolveUDPAddr("udp", "localhost:42069")

	if err != nil {
		log.Fatalf("failed to resolve: %s", err)
		return
	}

	conn, err := net.DialUDP("udp", nil, udp)

	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		m, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("failed to read %s", err)
		}

		conn.Write([]byte(m))

	}

}
