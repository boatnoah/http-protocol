package response

import (
	"fmt"
	"io"

	"github.com/httpfromtcp/internal/headers"
)

const httpVersion = "1.1"

type StatusCode int

const (
	Ok                  StatusCode = 200
	BadRequest          StatusCode = 400
	InternalServerError StatusCode = 500
)

type ReasonPhrase string

const (
	ReasonOk                  ReasonPhrase = "OK"
	ReasonBadRequest          ReasonPhrase = "Bad Request"
	ReasonInternalServerError ReasonPhrase = "Internal Server Error"
)

type Writer struct {
	dst         io.Writer
	statusCode  StatusCode
	hdrs        headers.Headers
	writerState int
}

func NewWriter(dst io.Writer) *Writer {
	return &Writer{dst: dst, writerState: 0}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.writerState != 0 {
		return fmt.Errorf("WriteStatusLine called in wrong order (state=%d)", w.writerState)
	}
	w.statusCode = statusCode
	w.writerState = 1
	return nil
}

func (w *Writer) WriteHeaders(h headers.Headers) error {
	if w.writerState != 1 {
		return fmt.Errorf("WriteHeaders called in wrong order (state=%d)", w.writerState)
	}
	c := headers.NewHeaders()
	for k, v := range h {
		c.Set(k, v)
	}
	w.hdrs = c
	w.writerState = 2
	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.writerState != 2 {
		return 0, fmt.Errorf("WriteBody called in wrong order (state=%d)", w.writerState)
	}
	if w.hdrs == nil {
		w.hdrs = headers.NewHeaders()
	}

	if _, ok := w.hdrs["connection"]; !ok {
		w.hdrs.Set("connection", "close")
	}
	w.hdrs.Set("content-length", fmt.Sprintf("%d", len(p)))

	if err := WriteStatusLine(w.dst, w.statusCode); err != nil {
		return 0, err
	}
	if err := WriteHeaders(w.dst, w.hdrs); err != nil {
		return 0, err
	}
	n, err := w.dst.Write(p)
	if err == nil {
		w.writerState = 3
	}
	return n, err
}

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	format := "HTTP/%s %d %s\r\n"
	switch statusCode {
	case Ok:
		_, err := io.WriteString(w, fmt.Sprintf(format, httpVersion, statusCode, ReasonOk))
		return err
	case BadRequest:
		_, err := io.WriteString(w, fmt.Sprintf(format, httpVersion, statusCode, ReasonBadRequest))
		return err
	case InternalServerError:
		_, err := io.WriteString(w, fmt.Sprintf(format, httpVersion, statusCode, ReasonInternalServerError))
		return err
	default:
		return fmt.Errorf("Unrecognized status code: %d", statusCode)
	}
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h.Set("content-length", fmt.Sprintf("%d", contentLen))
	h.Set("connection", "close")
	h.Set("content-type", "text/plain")
	return h
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for key, val := range headers {
		if _, err := io.WriteString(w, fmt.Sprintf("%s: %s\r\n", key, val)); err != nil {
			return err
		}
	}
	_, _ = io.WriteString(w, "\r\n")
	return nil
}
