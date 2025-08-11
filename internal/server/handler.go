package server

import (
	"io"

	"github.com/MrBhop/httpfromtcp/internal/request"
	"github.com/MrBhop/httpfromtcp/internal/response"
)

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

func (h *HandlerError) Write(w io.Writer) {
	headers := response.GetDefaultHeaders(len(h.Message))
	response.WriteStatusLine(w, h.StatusCode)
	response.WriteHeaders(w, headers)
	w.Write([]byte(h.Message))
}
