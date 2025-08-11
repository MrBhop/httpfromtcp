package server

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/MrBhop/httpfromtcp/internal/request"
	"github.com/MrBhop/httpfromtcp/internal/response"
)

type Server struct {
	closed atomic.Bool
	listener net.Listener
	handler Handler
}

func Serve(port int, handlerFunc Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	s := &Server{
		listener: listener,
		handler: handlerFunc,
	}
	go s.listen()
	return s, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			log.Printf("Error accepting connection: %s\n", err)
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	request, err := request.RequestFromReader(conn)
	if err != nil {
		handlerError := &HandlerError{
			StatusCode: response.StatusBadRequest,
			Message: err.Error(),
		}
		handlerError.Write(conn)
		return
	}
	outputBuffer := &bytes.Buffer{}

	if err := s.handler(outputBuffer, request); err != nil {
		err.Write(conn)
		return
	}

	headers := response.GetDefaultHeaders(outputBuffer.Len())
	if err := response.WriteStatusLine(conn, response.StatusOK); err != nil {
		log.Printf("Error writing status line: %s", err)
	}
	if err := response.WriteHeaders(conn, headers); err != nil {
		log.Printf("Error writing response headers: %s", err)
	}
	if _, err := conn.Write(outputBuffer.Bytes()); err != nil {
		log.Printf("Error writing response body: %s", err)
	}
}
