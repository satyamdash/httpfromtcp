package server

import (
	"fmt"
	"net"
	"strconv"

	"github.com/satyamdash/httpfromtcp/internal/request"
)

type ServerState int

type Bool struct {
	// contains filtered or unexported fields
}

const (
	Listening ServerState = iota
	Closed
)

type Server struct {
	listener net.Listener
}

func Serve(port int) (*Server, error) {
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}
	return &Server{listener: ln}, nil
}

func (s *Server) Close() error {
	return s.listener.Close()
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
	_, err := request.RequestFromReader(conn)

	if err != nil {
		fmt.Println(err)
	}
	resp := "HTTP/1.1 200 OK\r\n" +
		"Content-Type: text/plain\r\n" +
		"Content-Length: 13\r\n" +
		"\r\n" +
		"Hello World!"
	if _, err := conn.Write([]byte(resp)); err != nil {
		fmt.Println(err)
	}

	// fmt.Println("Request line:")
	// fmt.Printf("- Method: %s\n", r.RequestLine.Method)
	// fmt.Printf("- Target: %s\n", r.RequestLine.RequestTarget)
	// fmt.Printf("- Version: %s\n", r.RequestLine.HttpVersion)

	// fmt.Println("Headers:")

	// for key, val := range r.Headers {
	// 	fmt.Printf("- %s: %s\n", key, val)
	// }

	// fmt.Println("Body:")
	// fmt.Printf("%s\n", r.Body)

}
