package main

import (
	"log/slog"
	"net"

	"github.com/theaaronruss/go-http-server/internal"
)

func handleRequest(conn net.Conn) {
	defer conn.Close()
	slog.Info("Handling request for client")
	request, err := request.ReadRequest(conn)
	if err != nil {
		slog.Error(err.Error())
	}
	slog.Info("Parsed request", "request", request)
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		slog.Error("Failed to create listener: " + err.Error())
		return
	}
	defer listener.Close()
	slog.Info("Server started (" + listener.Addr().String() + ")")
	for {
		conn, err := listener.Accept()
		if err != nil {
			slog.Error("Failed to accept new connection: " + err.Error())
		}
		slog.Info("Accepted new connection from " + conn.RemoteAddr().String())
		go handleRequest(conn)
		if err != nil {
			slog.Error("Failed to send message to client: " + err.Error())
		}
	}
}
