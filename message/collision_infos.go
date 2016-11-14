package message

import (
	"barrage-server/ball"
	b "barrage-server/base"
	"barrage-server/libs/bufbo"
	"bytes"
	"fmt"
)

const (
	collisionInfoSize = 18
)

// CollisionInfo hold information about the collision between A and B.
type CollisionInfo struct {
	IDs     []b.FullBallID
	Damages []b.Damage
	States  []ball.State
}

// Size ...
func (ci *CollisionInfo) Size() int {
	return collisionInfoSize
}

// MarshalBinary ...
func (ci *CollisionInfo) MarshalBinary() ([]byte, error) {
	var buffer bytes.Buffer
	bw := bufbo.NewBEBufWriter(&buffer)

	// full ball id
	bw.PutUint32(uint32(ci.IDs[0].UID))
	bw.PutUint16(uint16(ci.IDs[0].ID))
	bw.PutUint32(uint32(ci.IDs[1].UID))
	bw.PutUint16(uint16(ci.IDs[1].ID))

	// damage
	bw.PutUint8(uint8(ci.Damages[0]))
	bw.PutUint8(uint8(ci.Damages[1]))

	// state
	isAlive, isKilled, err := ball.AnalyseStateToBytes(ci.States[0])
	if err != nil {
		return nil, fmt.Errorf(
			"%v the state is %d, while marshaling CollisionInfo-A(%v & %v).\n",
			err, ci.States[0], ci.IDs[0], ci.IDs[1])
	}
	bw.PutUint8(isAlive)
	bw.PutUint8(isKilled)
	isAlive, isKilled, err = ball.AnalyseStateToBytes(ci.States[1])
	if err != nil {
		return nil, fmt.Errorf(
			"%v the state is %d, while marshaling CollisionInfo-B(%v & %v).\n",
			err, ci.States[1], ci.IDs[0], ci.IDs[1])
	}
	bw.PutUint8(isAlive)
	bw.PutUint8(isKilled)
	// 28 bytes

	return buffer.Bytes(), nil
}

// UnmarshalBinary ...
func (ci *CollisionInfo) UnmarshalBinary(data []byte) error {
	ci.IDs = make([]b.FullBallID, 2)
	ci.Damages = make([]b.Damage, 2)
	ci.States = make([]ball.State, 2)

	br := bufbo.NewBEBytesReader(data)
	ci.IDs[0].UID = b.UserID(br.Uint32())
	ci.IDs[0].ID = b.BallID(br.Uint16())
	ci.IDs[1].UID = b.UserID(br.Uint32())
	ci.IDs[1].ID = b.BallID(br.Uint16())

	ci.Damages[0] = b.Damage(br.Uint8())
	ci.Damages[1] = b.Damage(br.Uint8())

	isAlive, isKilled := br.Uint8(), br.Uint8()
	state, err := ball.AnalyseBytesToState(isAlive, isKilled)
	if err != nil {
		return fmt.Errorf(
			"%v source bytes is isAlive: %d, isKilled: %d, while unmarshaling CollisionInfo-A(%v %v).\n",
			err, isAlive, isKilled, ci.IDs[0], ci.IDs[1])
	}
	ci.States[0] = state

	isAlive, isKilled = br.Uint8(), br.Uint8()
	state, err = ball.AnalyseBytesToState(isAlive, isKilled)
	if err != nil {
		return fmt.Errorf(
			"%v source bytes is isAlive: %d, isKilled: %d, while unmarshaling CollisionInfo-B(%v %v).\n",
			err, isAlive, isKilled, ci.IDs[0], ci.IDs[1])
	}
	ci.States[1] = state

	return nil
}

// CollisionsInfo is used for collision informations transimission.
type CollisionsInfo struct {
	length         uint32
	CollisionInfos []*CollisionInfo
}

// Length return length
func (csi *CollisionsInfo) Length() int {
	return int(csi.length)
}

// Item return item of CollisionInfos.
func (csi *CollisionsInfo) Item(index int) b.CommunicationData {
	return csi.CollisionInfos[index]
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
	csi.CollisionInfos = make([]*CollisionInfo, length)
	for i := range csi.CollisionInfos {
		csi.CollisionInfos[i] = new(CollisionInfo)
	}
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
