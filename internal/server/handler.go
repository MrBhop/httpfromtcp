package server

import (
	"github.com/MrBhop/httpfromtcp/internal/request"
	"github.com/MrBhop/httpfromtcp/internal/response"
)

type Handler func(w *response.Writer, req *request.Request)

func WriteConnectionError(w *response.Writer, message string) {
	headers := response.GetDefaultHeaders(len(message))
	w.WriteStatusLine(response.StatusBadRequest)
	w.WriteHeaders(headers)
	w.WriteBody([]byte(message))
}
