package main

import (
	"barrage-server/socket"
)

func main() {
	socket.ListenAndServer("2334", "/ws")
}
