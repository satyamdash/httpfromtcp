package request

import (
	"fmt"
	"io"
	"strings"
)

type StateLine int

const (
	Initialized StateLine = iota
	Done
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
	StateLine
}

const separator = "\r\n"
const HTTPprefix = "HTTP/"
const bufferSize = 8

func parseRequestLine(str string) (*RequestLine, int, error) {
	idx := strings.Index(str, "\r\n")
	if idx == -1 {
		return nil, 0, nil
	}
	request_line := str[:idx]
	request_info := strings.Fields(request_line)
	if len(request_info) < 3 {
		return nil, -1, fmt.Errorf("Not enough info")
	}

	var r RequestLine
	r.Method = request_info[0]
	r.RequestTarget = request_info[1]
	r.HttpVersion, _ = strings.CutPrefix(request_info[2], HTTPprefix)
	return &r, idx + len("\r\n"), nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize)
	var readtoindex = 0
	req := &Request{
		StateLine: Initialized,
	}

	for req.StateLine != Done {

		if readtoindex == len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		n, err := reader.Read(buf[readtoindex:])
		if err == io.EOF {
			req.StateLine = Done
			break
		}
		if err != nil {
			return nil, err
		}
		readtoindex += n

		consumed, err := req.Parse(buf[:readtoindex])
		if err != nil {
			return nil, err
		}
		if consumed > 0 {
			copy(buf, buf[consumed:readtoindex])
			readtoindex -= consumed
		}
	}
	return req, nil
}

func (r *Request) Parse(data []byte) (int, error) {
	switch r.StateLine {
	case Initialized:
		reqLine, byteLen, err := parseRequestLine(string(data))
		if err != nil {
			return 0, err
		}

		// Not enough data yet (no \r\n found)
		if reqLine == nil {
			return 0, nil
		}

		// Successfully parsed request line
		r.RequestLine = *reqLine
		r.StateLine = Done
		return byteLen, nil

	case Done:
		return 0, fmt.Errorf("error: trying to read data in a done state")

	default:
		return 0, fmt.Errorf("error: unknown state")
	}
}
