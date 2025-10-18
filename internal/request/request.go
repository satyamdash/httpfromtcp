package request

import (
	"fmt"
	"io"
	"strings"

	"github.com/satyamdash/httpfromtcp/internal/headers"
)

type StateLine int

const (
	Initialized StateLine = iota
	Done
	RequestStateParsingHeaders
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	StateLine
}

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
		return nil, -1, fmt.Errorf("not enough info")
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
		Headers:   headers.NewHeaders(),
	}

	for req.StateLine != Done {

		if readtoindex == len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		n, err := reader.Read(buf[readtoindex:])
		if n > 0 {
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

		if err == io.EOF {
			if req.StateLine != Done {
				return nil, fmt.Errorf("missing CRLF CRLF terminator at end of headers")
			}
			break
		}
		if err != nil {
			return nil, err
		}

	}
	return req, nil
}

func (r *Request) Parse(data []byte) (int, error) {
	totalBytesParsed := 0
	for r.StateLine != Done {
		n, err := r.ParseSingle(data[totalBytesParsed:])
		if err != nil {
			return 0, err
		}
		if n == 0 {
			break
		}
		totalBytesParsed += n
	}

	return totalBytesParsed, nil

}

func (r *Request) ParseSingle(data []byte) (int, error) {
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
		r.StateLine = RequestStateParsingHeaders
		return byteLen, nil

	case RequestStateParsingHeaders:
		n, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if done {
			r.StateLine = Done
		}
		return n, err

	case Done:
		return 0, fmt.Errorf("error: trying to read data in a done state")

	default:
		return 0, fmt.Errorf("error: unknown state")
	}
}
