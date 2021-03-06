package base

import (
	"barrage-server/libs/log"
	"errors"
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
type Damage uint8

// FullBallID is the full id of ball, it consists of UserID and BallID
type FullBallID struct {
	UID UserID
	ID  BallID
}

// Log is the logger of whole sys, it print to stdout and stderr in development env.
var Log log.Logger

var (
	// errors

	// ErrServerError defines all server error.
	ErrServerError = errors.New("Server error.")
)

const (
	// SysID is the user id for all ball and info created by server.
	SysID = 0
)

const (
	// environment

	// Development environment
	Development = iota
	// Testing environment
	Testing
	// Production environment
	Production
)

// init
func init() {
	Log = log.NewStdLogger(log.InfoLevel)
}
