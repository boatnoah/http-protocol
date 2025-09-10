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
