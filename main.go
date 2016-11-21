package main

import (
	b "barrage-server/base"
	r "barrage-server/room"
	"barrage-server/socket"
	"flag"
)

var env string

// TODO: write config package.
func init() {
	const (
		defaultEnv = "dev"
		usage      = "set running environment[dev, pro, test]"
	)
	flag.StringVar(&env, "env", defaultEnv, usage)
	flag.StringVar(&env, "e", defaultEnv, usage+" (shorthand)")
}

func main() {
	flag.Parse()

	// set running time environment
	switch env {
	case "test":
		b.RunningEnv = b.Testing
	case "pro":
		b.RunningEnv = b.Production
	case "dev":
		fallthrough
	default:
		b.RunningEnv = b.Development
	}

	r.OpenGameHallAndRooms(b.OpenRoomIDs)

	path := "/test"
	if b.RunningEnv == b.Production {
		path = "/ws"
	}

	socket.ListenAndServer("2334", path)
}
