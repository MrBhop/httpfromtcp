package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

const serverAddress = "localhost:42069"
const network = "udp"

func main() {
	udpAddress, err := net.ResolveUDPAddr(network, serverAddress)
	if err != nil {
		log.Fatalf("Error resolving address: %s\n", err)
	}

	conn, err := net.DialUDP(network, nil, udpAddress)
	if err != nil {
		log.Fatalf("Error opening connection: %s\n", err)
	}
	defer conn.Close()

	r := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf(">")
		line, err := r.ReadString('\n')
		if err != nil {
			log.Fatalf("Error reading from StdIn: %s\n", err)
		}

		if _, err := conn.Write([]byte(line)); err != nil {
			log.Fatalf("Error sending message: %s", err)
		}
	}
}
