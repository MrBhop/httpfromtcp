package request

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"
)


type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion string
	RequestTarget string
	Method string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(request), "\r\n")
	if length := len(lines); length < 2 {
		return nil, fmt.Errorf("Request does not have enough lines. Expected > 2, got %d", length)
	}

	requestLine, err := parseRequestLine(lines[0])
	if err != nil {
		return nil, fmt.Errorf("Error parsing request-line: %w", err)
	}

	output := Request {
		RequestLine: *requestLine,
	}

	return &output, nil
}

func parseRequestLine(line string) (*RequestLine, error) {
	fields := strings.Fields(line)
	if length := len(fields); length != 3 {
		return nil, fmt.Errorf("Request lines has incorrect number of whitespace delimited fields. Expected 3, got %d", length)
	}

	output := RequestLine {
		RequestTarget: fields[1],
	}

	// validate Method.
	method := fields[0]
	if upperCaseMethod := strings.ToUpper(method); upperCaseMethod != method {
		return nil, fmt.Errorf("Method contains non uppercase characters.")
	}

	for _, r := range method {
		if !unicode.IsLetter(r) {
			return nil, fmt.Errorf("Method contains non alphabetic characters.")
		}
	}

	output.Method = method

	// validate HTTP-version.
	// general validation.
	versionString := []byte(fields[2])
	if length := len(versionString); length != 8 {
		return nil, fmt.Errorf("HTTP-version is malformed. Expected length == 8, got %d", length)
	}

	if versionStart := string(versionString[:5]); versionStart != "HTTP/" {
		return nil, fmt.Errorf("HTTP-version is malformed. Expected to start with HTTP/, got: '%s'", versionStart)
	}

	digit1 := string(versionString[5])
	if _, err := strconv.Atoi(digit1); err != nil {
		return nil, fmt.Errorf("HTTP-version is malformed. Expected 6th char to be digit, got: '%s'", digit1)
	}

	if dot := string(versionString[6]); dot != "." {
		return nil, fmt.Errorf("HTTP-version is malformed. Expected 7th chat to be '.', got: '%s'", dot)
	}
	
	digit2 := string(versionString[7])
	if _, err := strconv.Atoi(digit2); err != nil {
		return nil, fmt.Errorf("HTTP-version is malformed. Expected 8th char to be digit, got: '%s'", digit1)
	}

	output.HttpVersion = digit1 + "." + digit2

	// version specific validation.
	if !(output.HttpVersion == "1.1") {
		return nil, fmt.Errorf("This application only supports 1.1, got: %s", output.HttpVersion)
	}

	return &output, nil
}

