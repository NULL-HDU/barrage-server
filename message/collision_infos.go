package message

import (
	"barrage-server/ball"
	b "barrage-server/base"
	"barrage-server/libs/bufbo"
	"bytes"
	"fmt"
)

type fullBallID struct {
	uid b.UserID
	id  b.BallID
}

// collisionInfo hold information about the collision between A and B.
type collisionInfo struct {
	ballIDs []fullBallID
	damages []b.Damage
	states  []ball.State
}

// AInfo return the ballId, damage and state of A.
func (ci *collisionInfo) AInfo() (b.UserID, b.BallID, b.Damage, ball.State) {
	return ci.ballIDs[0].uid, ci.ballIDs[0].id, ci.damages[0], ci.states[0]
}

// BInfo return the ballId, damage and state of B.
func (ci *collisionInfo) BInfo() (b.UserID, b.BallID, b.Damage, ball.State) {
	return ci.ballIDs[1].uid, ci.ballIDs[1].id, ci.damages[1], ci.states[1]
}

// Size ...
func (ci *collisionInfo) Size() int {
	return 28
}

// MarshalBinary ...
func (ci *collisionInfo) MarshalBinary() ([]byte, error) {
	var buffer bytes.Buffer
	bw := bufbo.NewBEBufWriter(&buffer)

	// full ball id
	bw.PutUint64(uint64(ci.ballIDs[0].uid))
	bw.PutUint16(uint16(ci.ballIDs[0].id))
	bw.PutUint64(uint64(ci.ballIDs[1].uid))
	bw.PutUint16(uint16(ci.ballIDs[1].id))

	// damage
	bw.PutUint16(uint16(ci.damages[0]))
	bw.PutUint16(uint16(ci.damages[1]))

	// state
	isAlive, isKilled, err := ball.AnalyseStateToBytes(ci.states[0])
	if err != nil {
		return nil, fmt.Errorf(
			"%v the state is %d, while marshaling collisionInfo-A(%v & %v).\n",
			err, ci.states[0], ci.ballIDs[0], ci.ballIDs[1])
	}
	bw.PutUint8(isAlive)
	bw.PutUint8(isKilled)
	isAlive, isKilled, err = ball.AnalyseStateToBytes(ci.states[1])
	if err != nil {
		return nil, fmt.Errorf(
			"%v the state is %d, while marshaling collisionInfo-B(%v & %v).\n",
			err, ci.states[1], ci.ballIDs[0], ci.ballIDs[1])
	}
	bw.PutUint8(isAlive)
	bw.PutUint8(isKilled)
	// 28 bytes

	return buffer.Bytes(), nil
}

// UnmarshalBinary ...
func (ci *collisionInfo) UnmarshalBinary(data []byte) error {
	ci.ballIDs = make([]fullBallID, 2)
	ci.damages = make([]b.Damage, 2)
	ci.states = make([]ball.State, 2)

	br := bufbo.NewBEBytesReader(data)
	ci.ballIDs[0].uid = b.UserID(br.Uint64())
	ci.ballIDs[0].id = b.BallID(br.Uint16())
	ci.ballIDs[1].uid = b.UserID(br.Uint64())
	ci.ballIDs[1].id = b.BallID(br.Uint16())

	ci.damages[0] = b.Damage(br.Uint16())
	ci.damages[1] = b.Damage(br.Uint16())

	isAlive, isKilled := br.Uint8(), br.Uint8()
	state, err := ball.AnalyseBytesToState(isAlive, isKilled)
	if err != nil {
		return fmt.Errorf(
			"%v source bytes is isAlive: %d, isKilled: %d, while unmarshaling collisionInfo-A(%v %v).\n",
			err, isAlive, isKilled, ci.ballIDs[0], ci.ballIDs[1])
	}
	ci.states[0] = state

	isAlive, isKilled = br.Uint8(), br.Uint8()
	state, err = ball.AnalyseBytesToState(isAlive, isKilled)
	if err != nil {
		return fmt.Errorf(
			"%v source bytes is isAlive: %d, isKilled: %d, while unmarshaling collisionInfo-B(%v %v).\n",
			err, isAlive, isKilled, ci.ballIDs[0], ci.ballIDs[1])
	}
	ci.states[1] = state

	return nil
}

// CollisionsInfo is used for collision informations transimission.
type CollisionsInfo struct {
	length         uint32
	CollisionInfos []collisionInfo
}

// Length return length
func (csi *CollisionsInfo) Length() int {
	return int(csi.length)
}

// Item return item of CollisionInfos.
func (csi *CollisionsInfo) Item(index int) b.CommunicationData {
	return &csi.CollisionInfos[index]
}

// Size return the number of bytes after marshed
func (csi *CollisionsInfo) Size() int {
	sum := 4
	for _, v := range csi.CollisionInfos {
		sum += v.Size()
	}
	return sum
}

// NewItems init CollisionInfos
func (csi *CollisionsInfo) NewItems(length uint32) {
	csi.CollisionInfos = make([]collisionInfo, length)
	csi.length = length
}

// Crop crop CollisionInfos
func (csi *CollisionsInfo) Crop(length uint32) {
	if csi.length == length {
		return
	}
	csi.CollisionInfos = csi.CollisionInfos[:length]
	csi.length = length
}
