// Package ball provide the interface and entity of ball
package ball

import (
	b "barrage-server/base"
	"barrage-server/libs/bufbo"
	"bytes"
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
	Alive = State(iota + 1)
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
	ballBaseSize = 26
)

type hp uint8

type role uint8

type special uint16

type speed uint8

type attackDir uint16

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

	HP() hp
	Damage() b.Damage
	SetHP(hp)

	// Position() (x, y uint16)
	IsDisappear() bool
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
	speed     speed
	attackDir attackDir
	state     State
	location  location
}

// NewUserAirplane creates an airplane for user.
//
// hp, damage, speed and attackDir is generate automatically according to roleConf
// of roleConfTable[r].
func NewUserAirplane(c b.UserID, nickname string, r role, s special, x uint16, y uint16) (Ball, error) {
	// TODO: we need a role table. Analyze from json file,
	//       but now we just write hard.
	airPlaneRole, ok := roleConfTable[r]
	if !ok {
		return nil, fmt.Errorf(
			"%v the role id is %d, While create ball(%d %d) via NewUserAirplane.\n",
			errInvalidRole, r, c, 0)
	}

	return &ball{
		uid:       c,
		id:        0,
		nickname:  nickname,
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

// NewBall create a nil-value ball
func NewBall() Ball {
	return &ball{}
}

func (bl *ball) UID() b.UserID {
	return bl.uid
}

func (bl *ball) ID() b.BallID {
	return bl.id
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
	var buffer bytes.Buffer
	bw := bufbo.NewBEBufWriter(&buffer)

	//uid(userId) + ballId(ballId) + ballType(Uint8) + hp(Uint16) + damage(damage)+
	//role(Uint8) + special(Uint16) + speed(Uint8) + attackDir(Uint16) + alive(bool) +
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
	bw.PutUint8(uint8(bl.speed))
	bw.PutUint16(uint16(bl.attackDir))

	isAlive, isKilled, err := AnalyseStateToBytes(bl.state)
	if err != nil {
		return nil, fmt.Errorf(
			"%v the state is %d, while marshaling ball(%d %d).\n",
			err, bl.state, bl.uid, bl.id)
	}
	bw.PutUint8(isAlive)
	bw.PutUint8(isKilled)
	bw.PutUint16(bl.location.x)
	bw.PutUint16(bl.location.y)
	// 34 + bytes

	return buffer.Bytes(), nil
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
	bl.speed = speed(br.Uint8())
	bl.attackDir = attackDir(br.Uint16())

	isAlive, isKilled := br.Uint8(), br.Uint8()
	state, err := AnalyseBytesToState(isAlive, isKilled)
	if err != nil {
		return fmt.Errorf(
			"%v source bytes is isAlive: %d, isKilled: %d, while unmarshaling ball(%d %d).\n",
			err, isAlive, isKilled, bl.uid, bl.id)
	}
	bl.state = state
	bl.location.x = br.Uint16()
	bl.location.y = br.Uint16()

	return nil
}

// AnalyseStateToBytes analyse state to isKilled and isAlive.
func AnalyseStateToBytes(s State) (uint8, uint8, error) {
	switch s {
	case Alive:
		return 1, 0, nil
	case Dead:
		return 0, 1, nil
	case Disappear:
		return 0, 0, nil
	default:
		return 0, 0, errInvalidState
	}
}

// AnalyseBytesToState analyse isKilled and isAlive to state
func AnalyseBytesToState(a, k uint8) (State, error) {
	switch {
	case a == 1 && k == 0:
		return Alive, nil
	case a == 0 && k == 1:
		return Dead, nil
	case a == 0 && k == 0:
		return Disappear, nil
	default:
		return 0, errInvalidState
	}
}
