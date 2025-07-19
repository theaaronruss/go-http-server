package main

import (
	"fmt"
	"net"
	"os"

	"github.com/theaaronruss/go-http-server/internal/request"
	"github.com/theaaronruss/go-http-server/internal/response"
)

func handleRequest(conn net.Conn) {
	defer conn.Close()
	request, err := request.ReadRequest(conn)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	switch request.Target {
	case "/error":
		response := response.NewResponse(400, []byte("This is an example of a bad request error\n"))
		response.Write(conn)
	case "/author":
		response := response.NewResponse(200, []byte("Aaron Russell\n"))
		response.Write(conn)
	default:
		response := response.NewResponse(404, nil)
		response.Write(conn)
	}
	conn.Close()
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to create listener:", err.Error())
		return
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to accept new connection:", err.Error())
		}
		go handleRequest(conn)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to send message to client:", err.Error())
		}
	}
}
