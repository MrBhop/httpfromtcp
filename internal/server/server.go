package server

import (
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

	response := response.Writer{
		Connection: conn,
	}

	request, err := request.RequestFromReader(conn)
	if err != nil {
		WriteConnectionError(&response, err.Error())
	}

	s.handler(&response, request)
}
