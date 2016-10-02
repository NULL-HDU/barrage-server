// Package ball provide the interface and entity of ball
package ball

import (
	"barrage-server/base"
)

type ballType uint8

const (
	airPlane = ballType(iota)
	block
	bullet
	food
)

type hp uint16

type damage uint16

type role uint8

type specail uint16

type speed uint8

type attackDir float32

type location struct {
	x float32
	y float32
}

// ballID is consist of userID and id. id is a value from 1 - 2^16.
// After user creating a ball, id add to one. 0 is user s airplane.
type ballID struct {
	userID base.UserID
	id     uint16
}

// Ball defines what a ball object should provide for backend.
//
// Other data should write in the struct.
type Ball interface {
	base.CommunicationData

	ID() ballID

	HP() hp
	Damage() damage
	SetHP(hp)

	IsDisappear() bool
}

// Collision do collisoin check and damage calculation for a and b.
// This function will return the true collisionInfo.
// func Collision(a Ball, b Ball) CollisionInfo {

// }

// NewBall creates instance of Ball from the binary.
//
// This should be used to receive ball data from frontend.
func NewBall(b []byte) Ball {

}
