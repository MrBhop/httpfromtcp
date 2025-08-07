package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

const filePath string = "messages.txt"

func main() {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Error opening file: %s", err)
	}
	defer file.Close()

	for {
		output := make([]byte, 8)
		_, err := file.Read(output)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			log.Fatalf("Error reading from file: %s", err)
		}

		fmt.Printf("read: %s\n", output)
	}
}
