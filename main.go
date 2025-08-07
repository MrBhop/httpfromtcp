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

	line := ""
	printAndResetLine := func() {
		fmt.Println("read:", line)
		line = ""
	}

	for {
		readBytes := make([]byte, 8)
		n, err := file.Read(readBytes)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			log.Fatalf("Error reading from file: %s", err)
		}

		readContentString := string(readBytes[:n])
		parts := strings.Split(readContentString, "\n")
		for i, p := range parts {
			line += p
			if i < len(parts) - 1 {
				printAndResetLine()
			}
		}
	}

	if line != "" {
		printAndResetLine()
	}
}
