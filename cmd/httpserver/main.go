package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/MrBhop/httpfromtcp/internal/request"
	"github.com/MrBhop/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	handlerFunc := func(w io.Writer, request *request.Request) *server.HandlerError {
		switch request.RequestLine.RequestTarget {
		case "/yourproblem":
			return &server.HandlerError{
				StatusCode: 400,
				Message: "Your problem is not my problem\n",
			}
		case "/myproblem":
			return &server.HandlerError{
				StatusCode: 500,
				Message: "Woopsie, my bad\n",
			}
		}
		w.Write([]byte("All good, frfr\n"))
		return nil
	}
	server, err := server.Serve(port, handlerFunc)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
