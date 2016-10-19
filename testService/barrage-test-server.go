// Package provide a test service for socket of frontend to test game protocal.
package main

import (
	"barrage-server/libs/log"
	"barrage-server/testService/data"
	"golang.org/x/net/websocket"
	"io"
	"net/http"
)

var logger log.Logger

func init() {
	logger = log.NewStdLogger(log.InfoLevel)
}

// echo response a test binary to client whenever client connect to server
func echo(ws *websocket.Conn) {
	logger.Infof("Connect from %v", ws.RemoteAddr())

	if err := websocket.Message.Send(ws, data.RandomUserID()); err != nil {
		logger.Errorf("Can't send: %s", err)
	}

	var cache []byte
	for {
		if err := websocket.Message.Receive(ws, &cache); err != nil {
			if err != io.EOF {
				logger.Errorf("Error: %s", err)
			}
			break
		}
		logger.Infof("Receive: % x\n", cache)

		if err := websocket.Message.Send(ws, cache); err != nil {
			logger.Errorf("Can't send: %s", err)
			break
		}
	}

	logger.Infof("Close Connect from %v", ws.RemoteAddr())
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
