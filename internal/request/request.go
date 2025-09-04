package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/httpfromtcp/internal/headers"
)

type parserState int

const (
	initialState parserState = iota
	parsingHeader
	parsingBody
	doneState
)

const cl = "Content-Length"

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	Body        []byte
	State       parserState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const crlf = "\r\n"

func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0

	for r.State != doneState {
		n, err := r.parseSingle(data[totalBytesParsed:])
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

func (r *Request) parseSingle(data []byte) (int, error) {

	switch r.State {

	case initialState:
		rl, consumed, err := parseRequestLine(data)

		if err != nil {
			return 0, nil
		}

		if consumed == 0 {
			return 0, nil
		}

		r.RequestLine = *rl

		r.State = parsingHeader

		return consumed, nil

	case parsingHeader:
		consumed, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, nil
		}

		if consumed == 0 {
			return 0, nil
		}

		if done {
			r.State = parsingBody
		}
		return consumed, nil

	case parsingBody:
		val, err := r.Headers.Get(cl)

		if err != nil {
			r.State = doneState
			return 0, nil
		}

		n, err := strconv.Atoi(val)
		if n <= 0 {
			r.State = doneState
			return 0, nil
		}

		if len(data) > 0 {
			r.Body = append(r.Body, data...)
		}

		if len(r.Body) > n {
			return 0, fmt.Errorf("body larger than Content-Length (expected %d, got %d+)", n, len(r.Body))
		}

		if len(r.Body) == n {
			r.State = doneState
		}

		return len(data), nil

	case doneState:
		return 0, nil
	}

	return 0, nil

}

func RequestFromReader(reader io.Reader) (*Request, error) {
	r := &Request{
		State:   initialState,
		Headers: headers.NewHeaders(),
	}
	buf := make([]byte, 8)
	readToIndex := 0

	for r.State != doneState {
		if readToIndex >= len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		numBytesRead, err := reader.Read(buf[readToIndex:])

		if err != nil {
			if errors.Is(err, io.EOF) {
				if r.State != doneState {
					return nil, fmt.Errorf("incomplete request (EOF before end of headers)")
				}
				break
			}
			return nil, err
		}
		readToIndex += numBytesRead

		numBytesParsed, err := r.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}
		copy(buf, buf[numBytesParsed:])
		readToIndex -= numBytesParsed

	}

	return r, nil
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return nil, 0, nil
	}
	requestLineText := string(data[:idx])
	requestLine, err := requestLineFromString(requestLineText)
	if err != nil {
		return nil, 0, err
	}
	return requestLine, idx + len([]byte(crlf)), nil
}

func requestLineFromString(str string) (*RequestLine, error) {
	parts := strings.Split(str, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("poorly formatted request-line: %s", str)
	}

	method := parts[0]
	for _, c := range method {
		if c < 'A' || c > 'Z' {
			return nil, fmt.Errorf("invalid method: %s", method)
		}
	}

	requestTarget := parts[1]

	versionParts := strings.Split(parts[2], "/")
	if len(versionParts) != 2 {
		return nil, fmt.Errorf("malformed start-line: %s", str)
	}

	httpPart := versionParts[0]
	if httpPart != "HTTP" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", httpPart)
	}
	version := versionParts[1]
	if version != "1.1" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", version)
	}

	return &RequestLine{
		Method:        method,
		RequestTarget: requestTarget,
		HttpVersion:   versionParts[1],
	}, nil
}
