package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/satyamdash/httpfromtcp/internal/headers"
	"github.com/satyamdash/httpfromtcp/internal/request"
	"github.com/satyamdash/httpfromtcp/internal/response"
	"github.com/satyamdash/httpfromtcp/internal/server"
)

const port = 42069

func handler(w *response.Writer, req *request.Request) {
	if req.RequestLine.Method == "GET" {
		if req.RequestLine.RequestTarget == "/video" {
			h := headers.NewHeaders()

			bytes, err := os.ReadFile("assets/vim.mp4")
			if err != nil {
				fmt.Println("error reading file")
			}
			w.WriteStatusLine(response.StatusOK)
			h["Content-Type"] = "video/mp4"
			h["Content-Length"] = fmt.Sprintf("%d", len(bytes))
			w.WriteHeaders(h)
			w.WriteBody(bytes)
		}
	}
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
		if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/") {
			route := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")

			resp, err := http.Get("https://httpbin.org/" + route)
			if err != nil {
				fmt.Println(err)
				w.WriteStatusLine(response.StatusServerError)
				w.WriteBody([]byte("Failed to fetch from httpbin.org"))
				return
			}
			defer resp.Body.Close()
			w.WriteStatusLine(response.StatusOK)
			h := headers.NewHeaders()
			h["Content-Type"] = "text/html"
			h["Transfer-Encoding"] = "chunked"
			h["Trailer"] = "X-Content-Sha256, X-Content-Length"
			w.WriteHeaders(h)

			buf := make([]byte, 1024)
			hash_body := []byte("")
			for {
				n, err := resp.Body.Read(buf)
				if n > 0 {
					hash_body = append(hash_body, buf[:n]...)
					w.WriteChunkedBody(buf[:n])
				}

				if err == io.EOF {
					break
				}
				if err != nil {
					fmt.Println("Error reading httpbin response:", err)
					break
				}
				fmt.Println(n)
			}
			w.WriteChunkedBodyDone()
			sum := sha256.Sum256(hash_body)
			trailer := headers.NewHeaders()
			fmt.Printf("%x", sum)
			fmt.Printf("%d", len(hash_body))
			trailer["X-Content-Sha256"] = fmt.Sprintf("%x", sum)
			trailer["X-Content-Length"] = fmt.Sprintf("%d", len(hash_body))
			err = w.WriteTrailers(trailer)
			if err != nil {
				fmt.Println(err)
			}
		}
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
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	server.Listen()

	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
