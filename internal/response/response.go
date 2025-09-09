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
		statusLine := fmt.Sprintf(format, httpVersion, statusCode, ReasonOk)
		if _, err := io.WriteString(w, statusLine); err != nil {
			return err
		}
	case BadRequest:
		statusLine := fmt.Sprintf(format, httpVersion, statusCode, ReasonBadRequest)
		if _, err := io.WriteString(w, statusLine); err != nil {
			return err
		}
	case InternalServerError:
		statusLine := fmt.Sprintf(format, httpVersion, statusCode, ReasonInternalServerError)
		if _, err := io.WriteString(w, statusLine); err != nil {
			return err
		}
	default:
		return fmt.Errorf("Unrecognized status code: %d", statusCode)
	}
	return nil
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
		headerContent := fmt.Sprintf("%s: %s\r\n", key, val)
		_, err := io.WriteString(w, headerContent)
		if err != nil {
			return err
		}
	}
	io.WriteString(w, "\r\n")

	return nil

}
