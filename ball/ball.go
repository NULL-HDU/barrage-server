// Package ball provide the interface and entity of ball
package ball

import (
	b "barrage-server/base"
	"barrage-server/libs/bufbo"
	"errors"
	"fmt"
	"math"
)

var (
	// errInvalidState throw while the state of ball is illegal value.
	errInvalidState = errors.New("Invalid state of ball.")
	// errInvalidRole throw while the role of ball is not included in roleConfTable.
	errInvalidRole = errors.New("Invalid role of ball.")
)

// State represent the status of ball (alive, deed, disappear)
type State uint8

const (
	// Alive  ball alive
	// isAlive: true, willDisappear: false
	Alive = State(iota)
	// Dead  ball dead
	// isAlive: false, willDisappear: false
	Dead
	// Disappear  ball disappear
	// isAlive: *, willDisappear: true
	Disappear
)

// Type is the type of ball
type Type uint8

const (
	// AirPlane user's ball(plane)
	AirPlane = Type(iota)
	// Block block created randomly by server to hinder plane from moving.
	Block
	// Bullet bullet created by user plane
	Bullet
	// Food food created randomly by server to make plane power up.
	Food
)

const (
	ballBaseSize = 28
)

type hp uint8

type role uint8

type special uint16

type radius uint16

type attackDir float32

type location struct {
	x uint16
	y uint16
}

// Ball defines what a ball object should provide for backend.
//
// Other data should write in the struct.
type Ball interface {
	b.CommunicationData

	UID() b.UserID
	ID() b.BallID

	// HP() hp
	// Damage() b.Damage
	// SetHP(hp)

	// // Position() (x, y uint16)
	// IsDisappear() bool
}

type ball struct {
	uid       b.UserID
	id        b.BallID
	nickname  string
	bType     Type
	hp        hp
	damage    b.Damage
	role      role
	special   special
	radius    radius
	attackDir attackDir
	state     State
	location  location
}

// NewBallFromBytes creates ball from bytes,
//
// The caller of this function should deal with all err and panic for index out of range!
func NewBallFromBytes(b []byte) (Ball, error) {
	newBall := &ball{}
	if err := newBall.UnmarshalBinary(b); err != nil {
		return nil, err
	}

	return newBall, nil
}

// NewBall create a nil-value ball
func NewBall() Ball {
	return &ball{}
}

// NewBallWithSpecialID create a nil-value ball
func NewBallWithSpecialID(uid b.UserID, id b.BallID) Ball {
	return &ball{
		uid:   uid,
		id:    id,
		state: Alive,
	}
}

func (bl *ball) UID() b.UserID {
	return bl.uid
}

func (bl *ball) ID() b.BallID {
	return bl.id
}

func (bl *ball) Radius() radius {
	return bl.radius
}

func (bl *ball) HP() hp {
	return bl.hp
}

func (bl *ball) Damage() b.Damage {
	return bl.damage
}

func (bl *ball) SetHP(HP hp) {
	bl.hp = HP
}

// func (bl *ball) Position() (x, y uint16) {
// 	return bl.location.x, bl.location.y
// }

func (bl *ball) IsDisappear() bool {
	if bl.state == Disappear {
		return true
	}

	return false
}

func (bl *ball) Size() int {
	return ballBaseSize + len(bl.nickname)
}

func (bl *ball) MarshalBinary() ([]byte, error) {
	bs := make([]byte, bl.Size())
	bw := bufbo.NewBEBytesWriter(bs)

	//uid(userId) + ballId(ballId) + ballType(Uint8) + hp(Uint16) + damage(damage)+
	//role(Uint8) + special(Uint16) + radius(Uint8) + attackDir(Uint16) + alive(bool) +
	//isKilled(bool) + locationCurrent(location)
	bw.PutUint32(uint32(bl.uid))
	bw.PutUint32(uint32(bl.uid))
	bw.PutUint16(uint16(bl.id))

	nicknameLen := len(bl.nickname)
	if nicknameLen > math.MaxUint8 {
		return nil, fmt.Errorf("Nickname is too long, hope 255, get %d.", nicknameLen)
	}
	bw.PutUint8(uint8(nicknameLen))
	bw.PutStr(bl.nickname)
	bw.PutUint8(uint8(bl.bType))
	bw.PutUint8(uint8(bl.hp))
	bw.PutUint8(uint8(bl.damage))
	bw.PutUint8(uint8(bl.role))
	bw.PutUint16(uint16(bl.special))
	bw.PutUint16(uint16(bl.radius))
	bw.PutFloat32(float32(bl.attackDir))
	bw.PutUint8(uint8(bl.state))
	bw.PutUint16(bl.location.x)
	bw.PutUint16(bl.location.y)
	// 28 + bytes

	return bs, nil
}

func (bl *ball) UnmarshalBinary(data []byte) error {
	br := bufbo.NewBEBytesReader(data)

	bl.uid = b.UserID(br.Uint32())
	bl.uid = b.UserID(br.Uint32())
	bl.id = b.BallID(br.Uint16())
	bl.nickname = br.Str(int(br.Uint8()))
	bl.bType = Type(br.Uint8())
	bl.hp = hp(br.Uint8())
	bl.damage = b.Damage(br.Uint8())
	bl.role = role(br.Uint8())
	bl.special = special(br.Uint16())
	bl.radius = radius(br.Uint16())
	bl.attackDir = attackDir(br.Float32())
	bl.state = State(br.Uint8())
	bl.location.x = br.Uint16()
	bl.location.y = br.Uint16()

	return nil
}
