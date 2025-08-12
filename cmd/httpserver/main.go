package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
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
	switch request.RequestLine.RequestTarget {
	case "/yourproblem":
		yourProblemHandler(w)
	case "/myproblem":
		myProblemHandler(w)
	default:
		if strings.HasPrefix(request.RequestLine.RequestTarget, "/httpbin/") {
			httpBinHandler(w, request.RequestLine.RequestTarget)
			return
		}
		okHandler(w)
	}
}

func httpBinHandler(w *response.Writer, target string) {
	nResponsesString := strings.TrimPrefix(target, "/httpbin/")

	resp, err := http.Get(fmt.Sprintf("https://httpbin.org/%s", nResponsesString))
	if err != nil {
		myProblemHandler(w)
		return
	}
	defer resp.Body.Close()

	headers := response.GetDefaultHeaders(0)
	headers.Remove("Content-Length")
	headers.Set("Content-Type", resp.Header.Get("Content-Type"))
	headers.Set("Transfer-Encoding", "chunked")

	w.WriteStatusLine(response.StatusOK)
	w.WriteHeaders(headers)

	for {
		buffer := make([]byte, 1024)
		n, err := resp.Body.Read(buffer)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			myProblemHandler(w)
			return
		}
		fmt.Printf("Read %d bytes\n", n)
		fmt.Printf("Bytes:\n%s\n", buffer)

		if _, err := w.WriteChunkedBody(buffer[:n]); err != nil {
			log.Println(err)
		}
	}
	if err := w.WriteChunkedBodyDone(); err != nil {
		log.Println(err)
	}
}

func yourProblemHandler(w *response.Writer) {
	statusCode := response.StatusBadRequest
	body := []byte(`<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`)
	basicHandler(w, body, statusCode)
}

func myProblemHandler(w *response.Writer) {
	statusCode := response.StatusInternalServerError
	body := []byte(`<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`)
	basicHandler(w, body, statusCode)
}

func okHandler(w *response.Writer) {
	statusCode := response.StatusOK
	body := []byte(`<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`)
	basicHandler(w, body, statusCode)
}

func basicHandler(w *response.Writer, body []byte, statusCode response.StatusCode) {
	headers := response.GetDefaultHeaders(len(body))
	headers.Set("Content-Type", "text/html")

	w.WriteStatusLine(statusCode)
	w.WriteHeaders(headers)
	w.WriteBody(body)
}
