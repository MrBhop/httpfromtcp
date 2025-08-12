package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/MrBhop/httpfromtcp/internal/constants"
	"github.com/MrBhop/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	StatusOK StatusCode = 200
	StatusBadRequest StatusCode= 400
	StatusInternalServerError StatusCode = 500
)

func GetStatusLine(statusCode StatusCode) []byte {
	var statusLine string
	switch statusCode {
	case StatusOK:
		statusLine = "OK"
	case StatusBadRequest:
		statusLine = "Bad Request"
	case StatusInternalServerError:
		statusLine = "Internal Server Error"
	default:
		return []byte{}
	}

	return fmt.Appendf([]byte{}, "HTTP/1.1 %d %s%s", statusCode, statusLine, constants.CrLf)
}

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	_, err := w.Write(GetStatusLine(statusCode))
	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	headers := headers.NewHeaders()
	headers.Add("Content-Length", strconv.Itoa(contentLen))
	headers.Add("Connection", "close")
	headers.Add("Content-Type", "text/plain")
	return headers
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for k, v := range headers {
		_, err := w.Write([]byte(k + ": " + v + constants.CrLf))
		if err != nil {
			return err
		}
	}

	_, err := w.Write([]byte(constants.CrLf))
	return err
}
