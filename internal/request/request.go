package request

import (
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"
)

const readBufferLen = 8

type requestState int

const (
	requestStateStatusLine requestState = iota
	requestStateHeaders
	requestStateBody
)

type Request struct {
	Method        string
	RequestTarget string
	HttpVersion   string
	Body          []byte
	Headers       map[string]string
	state         requestState
	currentLine   []byte
}

func RequestFromReader(r io.Reader) (*Request, error) {
	if readBufferLen <= 0 {
		return nil, errors.New("Buffer for reading request must be at least 1 byte in length")
	}
	request := &Request{ state: requestStateStatusLine }
	request.Headers = make(map[string]string)
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
	lineEndIndex := strings.Index(string(request.currentLine), "\r\n")
	if lineEndIndex == -1 {
		return nil
	}
	line := string(request.currentLine[:lineEndIndex + 2])
	request.currentLine = request.currentLine[lineEndIndex + 2:]
	var err error
	switch request.state {
	case requestStateStatusLine:
		err = parseRequestLine(request, line)
	case requestStateHeaders:
		err = parseHeader(request, line)
	}
	if err != nil {
		return err
	}
	return nil
}

func parseRequestLine(request *Request, requestLine string) error {
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

func parseHeader(request *Request, headerLine string) error {
	if headerLine == "\r\n" {
		request.state = requestStateBody
		return nil
	}
	separatorIndex := strings.Index(headerLine, ":")
	if separatorIndex == -1 {
		return errors.New("Invalid header")
	}
	headerName := headerLine[:separatorIndex]
	headerValue := headerLine[separatorIndex + 1:]
	headerName = strings.ToLower(headerName)
	headerValue = strings.TrimSpace(headerValue)
	request.Headers[headerName] = headerValue
	return nil
}
