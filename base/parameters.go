package base

import (
	"time"
)

// RunningEnv is the current running environment
var RunningEnv = Development

// RoomMembersLimit limit member in every room.
var RoomMembersLimit = 8

// PlayGroundHeight height of virtual playground
var PlayGroundHeight = 2100

// PlayGroundWidth width of virtual playground
var PlayGroundWidth = 3000

// OpenRoomIDs is the id of opened rooms
var OpenRoomIDs = []RoomID{1}

// RoomBoardCastDuration the duration between two boardcast of the room
var RoomBoardCastDuration = time.Millisecond * 40
