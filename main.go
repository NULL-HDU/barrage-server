package main

import (
	b "barrage-server/base"
	"barrage-server/socket"
)

func main() {

	path := "/test"

	if b.RunningEnv == b.Production {
		path = "/ws"
	}

	socket.ListenAndServer("2334", path)
}
