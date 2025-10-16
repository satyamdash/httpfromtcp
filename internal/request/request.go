package request

import (
	"fmt"
	"io"
	"strings"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
}

const separator = "\r\n"

func (r *RequestLine) parseRequestLine(str string) error {
	parts := strings.Split(str, separator)
	if len(parts) == 0 || strings.TrimSpace(parts[0]) == "" {
		return fmt.Errorf("invalid request: empty string")
	}
	req_line := parts[0]
	request_info := strings.Fields(req_line)
	if len(request_info) < 3 {
		return fmt.Errorf("invalid request line: %s", r)
	}
	r.Method = request_info[0]
	r.RequestTarget = request_info[1]
	const httpPrefix = "HTTP/"
	if !strings.HasPrefix(request_info[2], httpPrefix) {
		return fmt.Errorf("invalid HTTP version format: %q", request_info[2])
	}
	r.HttpVersion = request_info[2][len(httpPrefix):]
	return nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	bytearr, err := io.ReadAll(reader)
	if err != nil {
		fmt.Println(err)
	}

	str := string(bytearr)
	var r RequestLine
	if err := r.parseRequestLine(str); err != nil {
		return nil, err
	}

	req := &Request{
		RequestLine: r,
	}

	return req, nil

}
