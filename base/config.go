package base

import (
	"time"
)

// RoomMembersLimit limit member in every room.
var RoomMembersLimit = 8

// PlayGroundHeight height of virtual playground
var PlayGroundHeight = 3000

// PlayGroundWidth width of virtual playground
var PlayGroundWidth = 2100

// DefaultRoomID is the id of default room
var DefaultRoomID RoomID = 1

// DefaultRoomBoardCastDuration the duration between two boardcast of the room
var DefaultRoomBoardCastDuration = time.Millisecond * 10
