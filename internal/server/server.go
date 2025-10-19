package server

import (
	"fmt"
	"net"
	"strconv"

	"github.com/satyamdash/httpfromtcp/internal/request"
	"github.com/satyamdash/httpfromtcp/internal/response"
)

type ServerState int

type Bool struct {
	// contains filtered or unexported fields
}

type HandlerError struct {
	StatusCode int
	Message    string
}

type Handler func(w *response.Writer, req *request.Request)

const (
	Listening ServerState = iota
	Closed
)

type Server struct {
	listener net.Listener
	handler  Handler
}

func Serve(port int, handler Handler) (*Server, error) {
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}
	return &Server{listener: ln,
		handler: handler}, nil
}

func (s *Server) Close() error {
	return s.listener.Close()
}

func (s *Server) WriteError(w *response.Writer, err *HandlerError) {
	w.WriteStatusLine(response.StatusCode(err.StatusCode))

	// Write headers
	headers := response.GetDefaultHeaders(len(err.Message))
	w.WriteHeaders(headers)

	// Write body
	w.WriteBody([]byte(err.Message))
}

func (s *Server) Listen() {
	fmt.Println("Listening on", s.listener.Addr())
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Println(err)
			break
		}
		go s.handle(conn)
		fmt.Println("New Connection estblished")
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	r, err := request.RequestFromReader(conn)
	if err != nil {
		fmt.Println("Failed to parse request:", err)
		return
	}
	rw := response.NewWriter(conn)
	s.handler(rw, r)
}
