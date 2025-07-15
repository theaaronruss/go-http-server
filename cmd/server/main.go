package main

import (
	"log/slog"
	"net"
)

func handleRequest(conn net.Conn) {
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
		// err = sendMessage(conn)
		go handleRequest(conn)
		if err != nil {
			slog.Error("Failed to send message to client: " + err.Error())
		}
		conn.Close()
		slog.Info("Connection to client closed")
	}
}

// func sendMessage(conn net.Conn) error {
// 	slog.Info("Sending message to client")
// 	message := "Hello, user!\n"
// 	n, err := conn.Write([]byte(message))
// 	if err != nil {
// 		return err
// 	}
// 	slog.Info("Wrote " + fmt.Sprint(n) + " bytes to client")
// 	return nil
// }
