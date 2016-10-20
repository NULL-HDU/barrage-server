package base

import (
	"barrage-server/libs/log"
)

// UserID is id of user
type UserID uint32

// RoomID room.
type RoomID uint32

// BallID is a value from 1 - 2^16. After user creating a ball, id add to one.
// 0 is user s airplane.
type BallID uint16

// ImageID is id of image
type ImageID uint8

// Damage is damage of ball.
type Damage uint16

// Log is the logger of whole sys, it print to stdout and stderr in development env.
var Log log.Logger

// init
func init() {
	Log = log.NewStdLogger(log.InfoLevel)
}
