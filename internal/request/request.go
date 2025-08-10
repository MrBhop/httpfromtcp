package request

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"
)

type parserState int

const (
	initialized parserState = iota
	done
)

type Request struct {
	state parserState
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion string
	RequestTarget string
	Method string
}

const (
	bufferLength int = 8
	CrLf string = "\r\n"
)

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := &Request{state: initialized}
	buffer := make([]byte, bufferLength)
	usedBufferLength := 0

	for request.state != done {
		if capacity := cap(buffer); usedBufferLength >= capacity {
			newBuffer := make([]byte, capacity * 2)
			copy(newBuffer, buffer)
			buffer = newBuffer
		}

		n, err := reader.Read(buffer[usedBufferLength:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil, fmt.Errorf("Invalid request format - no crlf found")
			}

			return nil, err
		}

		usedBufferLength += n
		n, err = request.parse(buffer[:usedBufferLength])
		if err != nil {
			return nil, err
		}
	}

	return request, nil
}

func (r *Request) parse(next []byte) (int, error) {
	switch r.state {
	case initialized:
		n := strings.Index(string(next), CrLf)
		if n == -1 {
			return 0, nil
		}

		n, requestLine, err := parseRequestLine(next)
		if err != nil {
			return 0, err
		}
		r.RequestLine = *requestLine
		r.state = done
		return n, nil
	case done:
		return 0, fmt.Errorf("Cannot parse in a done state")
	default:
		return 0, fmt.Errorf("Unknown state")
	}
}

func parseRequestLine(line []byte) (int, *RequestLine, error) {
	n := strings.Index(string(line), CrLf)
	if n == -1 {
		return 0, nil, nil
	}

	lineString := string(line[:n])
	fields := strings.Fields(lineString)
	if length := len(fields); length != 3 {
		return 0, nil, fmt.Errorf("Request lines has incorrect number of whitespace delimited fields. Expected 3, got %d", length)
	}

	output := RequestLine {
		RequestTarget: fields[1],
	}

	// validate Method.
	method := fields[0]
	if upperCaseMethod := strings.ToUpper(method); upperCaseMethod != method {
		return 0, nil, fmt.Errorf("Method contains non uppercase characters.")
	}

	for _, r := range method {
		if !unicode.IsLetter(r) {
			return 0, nil, fmt.Errorf("Method contains non alphabetic characters.")
		}
	}

	output.Method = method

	// validate HTTP-version.
	// general validation.
	versionString := []byte(fields[2])
	if length := len(versionString); length != 8 {
		return 0, nil, fmt.Errorf("HTTP-version is malformed. Expected length == 8, got %d", length)
	}

	if versionStart := string(versionString[:5]); versionStart != "HTTP/" {
		return 0, nil, fmt.Errorf("HTTP-version is malformed. Expected to start with HTTP/, got: '%s'", versionStart)
	}

	digit1 := string(versionString[5])
	if _, err := strconv.Atoi(digit1); err != nil {
		return 0, nil, fmt.Errorf("HTTP-version is malformed. Expected 6th char to be digit, got: '%s'", digit1)
	}

	if dot := string(versionString[6]); dot != "." {
		return 0, nil, fmt.Errorf("HTTP-version is malformed. Expected 7th chat to be '.', got: '%s'", dot)
	}
	
	digit2 := string(versionString[7])
	if _, err := strconv.Atoi(digit2); err != nil {
		return 0, nil, fmt.Errorf("HTTP-version is malformed. Expected 8th char to be digit, got: '%s'", digit1)
	}

	output.HttpVersion = digit1 + "." + digit2

	// version specific validation.
	if !(output.HttpVersion == "1.1") {
		return 0, nil, fmt.Errorf("This application only supports 1.1, got: %s", output.HttpVersion)
	}

	return n + 2, &output, nil
}
