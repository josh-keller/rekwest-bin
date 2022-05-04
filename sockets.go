package main

import (
	"fmt"

	socketio "github.com/googollee/go-socket.io"
)

func (srv *server) socketRoutes() {
	srv.sockets.OnConnect("/", func(sock socketio.Conn) error {
		sock.SetContext("")
		fmt.Println("connected:", sock.ID())
		return nil
	})

	srv.sockets.OnEvent("/", "hello", func(sock socketio.Conn, msg string) string {
		sock.SetContext(msg)
		fmt.Println("Message: ", msg)
		return "recv " + msg
	})
}
