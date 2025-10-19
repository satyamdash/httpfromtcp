package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/satyamdash/httpfromtcp/internal/headers"
	"github.com/satyamdash/httpfromtcp/internal/request"
	"github.com/satyamdash/httpfromtcp/internal/response"
	"github.com/satyamdash/httpfromtcp/internal/server"
)

const port = 42069

func handler(w *response.Writer, req *request.Request) {
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		w.WriteStatusLine(response.StatusBadRequest)
		body := []byte(`<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`)
		h := headers.NewHeaders()
		h["Content-Type"] = "text/html"
		h["Content-Length"] = fmt.Sprintf("%d", len(body))
		w.WriteHeaders(h)
		w.WriteBody(body)

	case "/myproblem":
		w.WriteStatusLine(response.StatusServerError)
		body := []byte(`<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`)
		h := headers.NewHeaders()
		h["Content-Type"] = "text/html"
		h["Content-Length"] = fmt.Sprintf("%d", len(body))
		w.WriteHeaders(h)
		w.WriteBody(body)
	default:
		w.WriteStatusLine(response.StatusOK)
		body := []byte(`<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`)
		h := headers.NewHeaders()
		h["Content-Type"] = "text/html"
		h["Content-Length"] = fmt.Sprintf("%d", len(body))
		w.WriteHeaders(h)
		w.WriteBody(body)

	}
}
func main() {
	server, err := server.Serve(port, handler)
	defer server.Close()
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	server.Listen()

	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
