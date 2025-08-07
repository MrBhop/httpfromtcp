package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const filePath string = "messages.txt"

func main() {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Error opening file: %s", err)
	}
	defer file.Close()

	ch := getLinesChannel(file)
	for line := range ch {
		fmt.Println("read:", line)
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
