// Package ball provide the interface and entity of ball
package ball

import (
	"barrage-server/base"
	"barrage-server/libs/bufbo"
	"bytes"
	"errors"
)

// for guard
type userID base.UserID
type ballID base.BallID
type damage base.Damage

var logger = base.Log

type ballType uint8
type ballState uint8

var (
	// ErrInvalidState throw while the state of ball is illegal value.
	ErrInvalidState = errors.New("Invalid state of ball.")
	// ErrInvalidRole throw while the role of ball is not included in roleConfTable.
	ErrInvalidRole = errors.New("Invalid role of ball.")
)

const (
	// Alive  ball alive
	// isAlive: true, willDisappear: false
	Alive = ballState(iota)
	// Dead  ball dead
	// isAlive: false, willDisappear: false
	Dead
	// Disappear  ball disappear
	// isAlive: *, willDisappear: true
	Disappear
)

const (
	// AirPlane user's ball(plane)
	AirPlane = ballType(iota)
	// Block block created randomly by server to hinder plane from moving.
	Block
	// Bullet bullet created by user plane
	Bullet
	// Food food created randomly by server to make plane power up.
	Food
)

type hp uint16

type role uint8

type special uint16

type speed uint8

type attackDir float32

type location struct {
	x float32
	y float32
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

type ball struct {
	camp      userID
	id        ballID
	bType     ballType
	hp        hp
	damage    damage
	role      role
	special   special
	speed     speed
	attackDir attackDir
	state     ballState
	location  location
}

// NewUserAirplane creates an airplane for user.
//
// hp, damage, speed and attackDir is generate automatically according to roleConf
// of roleConfTable[r].
func NewUserAirplane(c userID, r role, s special, x float32, y float32) (Ball, error) {
	// TODO: we need a role table. Analyze from json file,
	//       but now we just write hard.
	airPlaneRole, ok := roleConfTable[r]
	if !ok {
		logger.Errorf("%s the role id is %d.\n", ErrInvalidRole, r)
		return nil, ErrInvalidRole
	}

	return &ball{
		camp:      c,
		id:        0,
		bType:     AirPlane,
		hp:        airPlaneRole.hp,
		damage:    airPlaneRole.damage,
		role:      r,
		special:   s,
		speed:     airPlaneRole.speed,
		attackDir: airPlaneRole.attackDir,
		state:     Alive,
		location:  location{x, y},
	}, nil
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

func (b *ball) ID() ballID {
	return b.id
}

func (b *ball) HP() hp {
	return b.hp
}

func (b *ball) Damage() damage {
	return b.damage
}

func (b *ball) SetHP(HP hp) {
	b.hp = HP
}

func (b *ball) IsDisappear() bool {
	if b.state == Disappear {
		return true
	}

	return false
}

func (b *ball) MarshalBinary() (data []byte, err error) {
	var buffer bytes.Buffer
	bw := bufbo.NewBEBufWriter(&buffer)

	//camp(userId) + ballId(ballId) + ballType(Uint8) + hp(Uint16) + damage(damage)+
	//role(Uint8) + special(Uint16) + speed(Uint8) + attackDir(Float32) + alive(bool) +
	//isKilled(bool) + locationCurrent(location)
	bw.PutUint64(uint64(b.camp))
	bw.PutUint64(uint64(b.camp))
	bw.PutUint16(uint16(b.id))
	bw.PutUint8(uint8(b.bType))
	bw.PutUint16(uint16(b.hp))
	bw.PutUint16(uint16(b.damage))
	bw.PutUint8(uint8(b.role))
	bw.PutUint16(uint16(b.special))
	bw.PutUint8(uint8(b.speed))
	bw.PutFloat32(float32(b.attackDir))
	switch b.state {
	case Alive:
		bw.PutUint8(1)
		bw.PutUint8(0)
	case Dead:
		bw.PutUint8(0)
		bw.PutUint8(1)
	case Disappear:
		bw.PutUint8(0)
		bw.PutUint8(0)
	default:
		return nil, ErrInvalidState
	}
	bw.PutFloat32(b.location.x)
	bw.PutFloat32(b.location.y)
	// 41 bytes

	data = buffer.Bytes()
	return
}

func (b *ball) UnmarshalBinary(data []byte) error {
	br := bufbo.NewBEBytesReader(data)

	b.camp = userID(br.Uint64())
	b.camp = userID(br.Uint64())
	b.id = ballID(br.Uint16())
	b.bType = ballType(br.Uint8())
	b.hp = hp(br.Uint16())
	b.damage = damage(br.Uint16())
	b.role = role(br.Uint8())
	b.special = special(br.Uint16())
	b.speed = speed(br.Uint8())
	b.attackDir = attackDir(br.Float32())
	switch isAlive, isKilled := br.Uint8(), br.Uint8(); {
	case isAlive == 1 && isKilled == 0:
		b.state = Alive
	case isAlive == 0 && isKilled == 1:
		b.state = Dead
	case isAlive == 0 && isKilled == 0:
		b.state = Disappear
	default:
		return ErrInvalidState
	}
	b.location.x = br.Float32()
	b.location.y = br.Float32()

	return nil
}
