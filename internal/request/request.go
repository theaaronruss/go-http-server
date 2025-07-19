package request

import (
	"errors"
	"fmt"
	"io"
	"net"
	"slices"
	"strconv"
	"strings"
)

const initialBufferLen = 8
var httpMethods = []string { "GET", "HEAD", "POST", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE" }

type requestState int

const (
	requestStateStatusLine requestState = iota
	requestStateHeaders
	requestStateBody
	requestStateDone
)

type Request struct {
	Method        string
	Target string
	HttpVersion   string
	Headers       map[string]string
	Body          []byte

	state         requestState
}

func ReadRequest(conn net.Conn) (*Request, error) {
	request := &Request{
		Headers: make(map[string]string),
		state: requestStateStatusLine,
	}
	buffer := make([]byte, initialBufferLen)
	readBytes := 0
	parsedBytes := 0
	for request.state != requestStateDone {
		if readBytes >= len(buffer) {
			newBuffer := make([]byte, len(buffer) * 2)
			copy(newBuffer, buffer)
			buffer = newBuffer
		}
		n, err := conn.Read(buffer[readBytes:])
		if err == io.EOF {
			request.state = requestStateDone
			return request, nil
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read request: %w", err)
		}
		readBytes += n
		p := -1
		for p != 0 {
			p, err = request.parse(buffer[parsedBytes:readBytes])
			if err != nil {
				return nil, fmt.Errorf("failed to parse request: %w", err)
			}
			parsedBytes += p
		}
	}
	return request, nil
}

func (request *Request) parse(data []byte) (int, error) {
	switch request.state {
	case requestStateStatusLine:
		return request.parseRequestLine(data)
	case requestStateHeaders:
		return request.parseHeader(data)
	case requestStateBody:
		if len(request.Body) == 0 {
			contentLengthHeader, ok := request.Headers["content-length"]
			if !ok {
				request.state = requestStateDone
				return 0, nil
			}
			contentLength, err := strconv.Atoi(contentLengthHeader)
			if err != nil {
				request.state = requestStateDone
				return 0, nil
			}
			request.Body = make([]byte, 0, contentLength)
		}
		return request.parseBody(data), nil
	default:
		return 0, nil
	}
}

func (request *Request) parseRequestLine(data []byte) (int, error) {
	dataStr := string(data)
	lineEndIndex := strings.Index(dataStr, "\r\n")
	if lineEndIndex == -1 {
		return 0, nil
	}
	parts := strings.Split(dataStr[:lineEndIndex], " ")
	if len(parts) != 3 {
		return 0, errors.New("invalid request line")
	}
	method := strings.ToUpper(parts[0])
	if !slices.Contains(httpMethods, method) {
		return 0, errors.New("invalid http method")
	}
	target := parts[1]
	version := parts[2]
	request.Method = method
	request.Target = target
	request.HttpVersion = version
	request.state = requestStateHeaders
	return lineEndIndex + 2, nil
}

func (request *Request) parseHeader(data []byte) (int, error) {
	dataStr := string(data)
	lineEndIndex := strings.Index(dataStr, "\r\n")
	if lineEndIndex == -1 {
		return 0, nil
	}
	if lineEndIndex == 0 {
		request.state = requestStateBody
		return 2, nil
	}
	colonIndex := strings.Index(dataStr, ":")
	if colonIndex == -1 {
		return 0, errors.New("invalid header")
	}
	headerName := strings.ToLower(dataStr[:colonIndex])
	headerValue := strings.TrimSpace(dataStr[colonIndex + 1:lineEndIndex])
	request.Headers[headerName] = headerValue
	return lineEndIndex + 2, nil
}

func (request *Request) parseBody(data []byte) int {
	request.Body = append(request.Body, data...)
	if len(request.Body) == cap(request.Body) {
		request.state = requestStateDone
	}
	return len(data)
}
