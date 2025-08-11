package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/MrBhop/httpfromtcp/internal/request"
	"github.com/MrBhop/httpfromtcp/internal/response"
	"github.com/MrBhop/httpfromtcp/internal/server"
)

const port = 42069

func main() {
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

func handlerFunc(w *response.Writer, request *request.Request) {
	var statusCode response.StatusCode
	var body []byte

	switch request.RequestLine.RequestTarget {
	case "/yourproblem":
		statusCode = response.StatusBadRequest
		body = []byte(`<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`)
	case "/myproblem":
		statusCode = response.StatusInternalServerError
		body = []byte(`<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`)
	default:
		statusCode = response.StatusOK
		body = []byte(`<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`)
	}

	headers := response.GetDefaultHeaders(len(body))
	headers.Set("Content-Type", "text/html")

	w.WriteStatusLine(statusCode)
	w.WriteHeaders(headers)
	w.WriteBody(body)
}
