package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/satyamdash/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	StatusOK          StatusCode = 200
	StatusBadRequest  StatusCode = 400
	StatusServerError StatusCode = 500
)

type Writer struct {
	Write io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{Write: w}
}
func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	switch statusCode {
	case StatusOK:
		_, err := w.Write.Write([]byte("HTTP/1.1 200 OK\r\n"))
		return err
	case StatusBadRequest:
		_, err := w.Write.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
		return err
	case StatusServerError:
		_, err := w.Write.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
		return err
	default:
		_, err := w.Write.Write([]byte(""))
		return err
	}
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	// Write headers
	for k, v := range headers {
		if _, err := w.Write.Write([]byte(fmt.Sprintf("%s: %s\r\n", k, v))); err != nil {
			return err
		}
	}
	_, err := w.Write.Write([]byte("\r\n"))
	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	return w.Write.Write(p)
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := make(headers.Headers)
	h["Content-Length"] = strconv.Itoa(contentLen)
	h["Connection"] = "close"
	h["Content-Type"] = "text/plain"

	return h
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	sizeLine := fmt.Sprintf("%x\r\n", len(p))
	if _, err := w.Write.Write([]byte(sizeLine)); err != nil {
		return 0, err
	}

	// Write the actual chunk data
	if _, err := w.Write.Write(p); err != nil {
		return 0, err
	}

	// Write the trailing CRLF after the chunk
	if _, err := w.Write.Write([]byte("\r\n")); err != nil {
		return 0, err
	}

	return len(p), nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	return w.Write.Write([]byte("0\r\n"))
}

func (w *Writer) WriteTrailers(h headers.Headers) error {
	for k, v := range h {
		if _, err := w.Write.Write([]byte(fmt.Sprintf("%s: %s\r\n", k, v))); err != nil {
			return err
		}
	}
	_, err := w.Write.Write([]byte("\r\n"))
	return err
}
