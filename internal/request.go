package request

import (
	"fmt"
	"net"
)

const initialBufferLen = 8
var httpMethods = []string { "GET", "HEAD", "POST", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE" }

type requestState int

const (
	requestStateStatusLine requestState = iota
	requestStateHeaders
	requestStateDone
)

type Request struct {
	Method        string
	RequestTarget string
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
	// parsedBytes := 0
	for request.state != requestStateDone {
		if readBytes >= len(buffer) {
			newBuffer := make([]byte, len(buffer) * 2)
			copy(newBuffer, buffer)
			buffer = newBuffer
		}
		n, err := conn.Read(buffer[readBytes:])
		if err != nil {
			return nil, fmt.Errorf("failed to read request: %w", err)
		}
		readBytes += n
	}
	return nil, nil
}





// func RequestFromReader(conn net.Conn) (*Request, error) {
// 	if initialBufferLength <= 0 {
// 		return nil, errors.New("size of buffer for reading must be at least 1 byte")
// 	}
// 	request := &Request{
// 		Headers: make(map[string]string),
// 		state: requestStateStatusLine,
// 	}
// 	readBytes := 0
// 	parsedBytes := 0
// 	buffer := make([]byte, initialBufferLength)
// 	for request.state != requestStateDone {
// 		if readBytes >= len(buffer) {
// 			buffer = slices.Grow(buffer, len(buffer))
// 		}
// 		n, err := conn.Read(buffer[readBytes:])
// 		fmt.Println(string(buffer))
// 		fmt.Println(cap(buffer))
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to read request: %w", err)
// 		}
// 		readBytes += n
// 		// fmt.Printf("Start: %d, End: %d\n", parsedBytes, readBytes)
// 		p, err := request.parse(buffer[parsedBytes:readBytes])
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to parse request: %w", err)
// 		}
// 		parsedBytes += p
// 	}
// 	return request, nil
// }
//
// func (request *Request) parse(data []byte) (int, error) {
// 	fmt.Println(string(data))
// 	switch request.state {
// 	case requestStateStatusLine:
// 		return request.parseRequestLine(data)
// 	case requestStateHeaders:
// 		return request.parseHeader(data)
// 	default:
// 		return 0, errors.New("Request entered invalid state while parsing")
// 	}
// }
//
// func (request *Request) parseRequestLine(data []byte) (int, error) {
// 	dataStr := string(data)
// 	lineEndIndex := strings.Index(dataStr, "\r\n")
// 	if lineEndIndex == -1 {
// 		return 0, nil
// 	}
// 	parts := strings.Split(dataStr[:lineEndIndex], " ")
// 	if len(parts) != 3 {
// 		return 0, errors.New("invalid request line")
// 	}
// 	request.Method = strings.ToUpper(parts[0])
// 	if !slices.Contains(httpMethods, request.Method) {
// 		return 0, errors.New("invalid http method")
// 	}
// 	request.RequestTarget = parts[1]
// 	request.HttpVersion = parts[2]
// 	request.state = requestStateHeaders
// 	return lineEndIndex + 2, nil
// }
//
// func (request *Request) parseHeader(data []byte) (int, error) {
// 	dataStr := string(data)
// 	if dataStr == "\r\n" {
// 		request.state = requestStateDone
// 		return 2, nil
// 	}
// 	lineEndIndex := strings.Index(dataStr, "\r\n")
// 	if lineEndIndex == -1 {
// 		return 0, nil
// 	}
// 	separatorIndex := strings.Index(dataStr, ":")
// 	if separatorIndex == -1 {
// 		return 0, errors.New("invalid header")
// 	}
// 	name := strings.ToLower(dataStr[:separatorIndex])
// 	value := strings.ToLower(dataStr[separatorIndex + 1:])
// 	request.Headers[name] = value
// 	return lineEndIndex, nil
// }
