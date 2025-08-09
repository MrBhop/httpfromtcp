package main

import (
	"fmt"
	"log"
	"net"

	"github.com/MrBhop/httpfromtcp/internal/request"
)

const port = ":42069"

func main() {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Error listening to TCP traffic: %s", err)
	}
	defer listener.Close()

	fmt.Println("Listening to TCP traffic on", port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Error accepting connection: %s", err)
		}

		fmt.Println("Connection accepted!")


		request, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatalf("Error parsing request: %s", err)
		}

		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n", request.RequestLine.Method)
		fmt.Printf("- Target: %s\n", request.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", request.RequestLine.HttpVersion)


		fmt.Println("Connection to", conn.RemoteAddr(), "closed!")
	}
}
