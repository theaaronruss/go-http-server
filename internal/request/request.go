package request

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"slices"
	"strings"
)

const readBufferLen = 8

type requestState int

const (
	requestStateStatusLine requestState = iota
	requestStateHeaders
)

type Request struct {
	Method        string
	RequestTarget string
	HttpVersion   string
	Body          []byte
	state         requestState
	currentLine   []byte
}

func RequestFromReader(r io.Reader) (*Request, error) {
	if readBufferLen <= 0 {
		return nil, errors.New("Buffer for reading request must be at least 1 byte in length")
	}
	request := &Request{ state: requestStateStatusLine }
	buffer := make([]byte, readBufferLen)
	for {
		n, err := r.Read(buffer)
		err = parse(request, buffer[:n])
		if err != nil {
			return nil, fmt.Errorf("Failed to read request: %w", err)
		}
		if errors.Is(err, io.ErrUnexpectedEOF) {
			break
		}
	}
	return request, nil
}

func parse(request *Request, newData []byte) error {
	if len(newData) == 0 {
		return nil
	}
	request.currentLine = slices.Concat(request.currentLine, newData)
	requestLineEndIndex := strings.Index(string(request.currentLine), "\r\n")
	if requestLineEndIndex == -1 {
		return nil
	}
	switch request.state {
	case requestStateStatusLine:
		requestLineStr := string(request.currentLine[:requestLineEndIndex + 2])
		err := parseRequestLine(request, requestLineStr)
		if err != nil {
			return err
		}
		request.currentLine = request.currentLine[requestLineEndIndex + 2:]
	}
	return nil
}

func parseRequestLine(request *Request, requestLine string) error {
	slog.Info("Request line: " + requestLine)
	parts := strings.Split(requestLine, " ")
	if len(parts) != 3 {
		return errors.New("Invalid request line")
	}
	request.Method = parts[0]
	request.RequestTarget = parts[1]
	request.HttpVersion = parts[2]
	request.state = requestStateHeaders
	return nil
}
