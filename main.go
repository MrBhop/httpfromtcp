package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
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
		ch := getLinesChannel(conn)
		for line := range ch {
			fmt.Println(line)
		}

		fmt.Println("Connection to", conn.RemoteAddr(), "closed!")
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)

	go func() {
		defer f.Close()
		defer close(ch)

		line := ""
		sendAndResetLine := func() {
			ch <- line
			line = ""
		}

		for {
			buffer := make([]byte, 8)
			n, err := f.Read(buffer)
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				log.Fatalf("Error reading from file: %s", err)
			}

			readContentString := string(buffer[:n])
			parts := strings.Split(readContentString, "\n")
			for i, p := range parts {
				line += p
				if i < len(parts) - 1 {
					sendAndResetLine()
				}
			}
		}
		if line != "" {
			sendAndResetLine()
		}
	}()

	return ch
}
