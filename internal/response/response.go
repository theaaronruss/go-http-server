package response

import (
	"bytes"
	"fmt"
	"net"
)

type Response struct {
	HttpVersion       string
	Status            int
	StatusDescription string
	Headers           map[string]string
	Body              []byte
}

func NewResponse(status int, body []byte) *Response {
	response := &Response{
		HttpVersion: "HTTP/1.1",
		Status: status,
		Headers: make(map[string]string),
		Body: body,
	}
	return response
}

func (response *Response) Write(conn net.Conn) error {
	response.Headers["content-length"] = fmt.Sprint(len(response.Body))
	var buffer bytes.Buffer
	buffer.WriteString(response.HttpVersion)
	buffer.WriteString(" ")
	buffer.WriteString(fmt.Sprint(response.Status))
	buffer.WriteString(" \r\n")
	for name, value := range response.Headers {
		buffer.WriteString(name + ": " + value + "\r\n")
	}
	buffer.WriteString("\r\n")
	buffer.Write(response.Body)
	_, err := conn.Write(buffer.Bytes())
	return err
}
