package response

import (
	"fmt"
	"net"

	"github.com/MrBhop/httpfromtcp/internal/constants"
	"github.com/MrBhop/httpfromtcp/internal/headers"
)

type WriterState int

const (
	WriterStatusLine WriterState = 0
	WriterHeaders = iota
	WriterBody
)

type Writer struct {
	writerState WriterState
	Connection net.Conn
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.writerState != WriterStatusLine {
		return fmt.Errorf("Invalid operation in the current state")
	}
	if err := WriteStatusLine(w.Connection, statusCode); err != nil {
		return err
	}
	w.writerState = WriterHeaders
	return nil
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.writerState != WriterHeaders {
		return fmt.Errorf("Invalid operation in the current state")
	}
	if err := WriteHeaders(w.Connection, headers); err != nil {
		return err
	}
	w.writerState = WriterBody
	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.writerState != WriterBody {
		return 0, fmt.Errorf("Invalid operation in the current state")
	}
	return w.Connection.Write(p)
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	crlfBytes := []byte(constants.CrLf)
	crlfLength := len(crlfBytes)
	bodyLength := len(p) + crlfLength
	bodyLengthLine := fmt.Appendf(nil, "%x%s", bodyLength, crlfBytes)

	completeBodyChunk := make([]byte, 0, len(bodyLengthLine) + bodyLength)
	completeBodyChunk = append(completeBodyChunk, bodyLengthLine...)
	completeBodyChunk = append(completeBodyChunk, p...)
	completeBodyChunk = append(completeBodyChunk, crlfBytes...)
	return w.WriteBody(completeBodyChunk)
}

func (w *Writer) WriteChunkedBodyDone() error {
	_, err := w.WriteChunkedBody([]byte{})
	return err
}
