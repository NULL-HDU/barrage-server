package message

import (
	"barrage-server/base"
)

//for guard
type ballID base.BallID
type damage base.Damage

type ballState uint8

const (
	// isAlive: true, willDisappear: false
	alive = ballState(iota)
	// isAlive: false, willDisappear: false
	dead
	// isAlive: *, willDisappear: true
	disappear
)

// CollisionInfo is the information about the collision causing between two balls,
// frontend socket -> backend information dosen't contain damages, but backend ->
// frontend information must contain damages!
type CollisionInfo interface {
	base.CommunicationData

	BallCauseCollision() []ballID
	BallStateAfterCollision() []ballState

	SetBallDamages([]damage)
}

// After Playground receiving a collisionInfo, it get id of balls via BallCauseCollision, then
// process checking function which will return damages. Playground set damages to collisionInfo
// before sending it to room.
