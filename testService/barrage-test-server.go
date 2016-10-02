// Package provide a test service for socket of frontend to test game protocal.
package main

import (
	"barrage-server/libs/log"
	"barrage-server/testService/data"
	"golang.org/x/net/websocket"
	"net/http"
)

var logger log.Logger

func init() {
	logger = log.NewStdLogger("TestService", log.InfoLevel)
}

// echo response a test binary to client whenever client connect to server
func echo(ws *websocket.Conn) {
	logger.Infof("Connect from %v", ws.RemoteAddr())

	if err := websocket.Message.Send(ws, data.Reply); err != nil {
		logger.Warnf("Can't send: %s", err)
	}

	ws.Close()
}

func main() {
	// provide file server
	http.Handle("/", http.FileServer(http.Dir("./public")))
	// provide websocket server
	http.Handle("/ws", websocket.Handler(echo))
	port := "1234"

	logger.Infof("Service start, bind port: %v", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		logger.Fatalln("ListenAndServe:", err)
	}
}
