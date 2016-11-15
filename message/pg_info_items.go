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
	disappearInfoSize = 2
)

// DisappearsInfo is used for disappear informations transimission.
type DisappearsInfo struct {
	IDs []b.BallID
}

// Size ...
func (dsi *DisappearsInfo) Size() int {
	return 4 + len(dsi.IDs)*disappearInfoSize
}

// MarshalBinary ...
func (dsi *DisappearsInfo) MarshalBinary() ([]byte, error) {
	length := len(dsi.IDs)
	bs := make([]byte, length*disappearInfoSize+4)
	bw := bufbo.NewBEBytesWriter(bs)

	bw.PutUint32(uint32(length))
	for _, id := range dsi.IDs {
		bw.PutUint16(uint16(id))
	}

	return bs, nil
}

// UnmarshalBinary ...
func (dsi *DisappearsInfo) UnmarshalBinary(data []byte) error {
	br := bufbo.NewBEBytesReader(data)

	length := br.Uint32()
	dsi.IDs = make([]b.BallID, length)
	for i := uint32(0); i < length; i++ {
		dsi.IDs[i] = b.BallID(br.Uint16())
	}

	return nil
}

// BallsInfo is used for ball informations transimission.
type BallsInfo struct {
	BallInfos []ball.Ball
}

// Length return length
func (bsi *BallsInfo) Length() int {
	return len(bsi.BallInfos)
}

// Item return item of BallInfos.
func (bsi *BallsInfo) Item(index int) b.CommunicationData {
	return bsi.BallInfos[index]
}

// Size return the number of bytes after marshed
func (bsi *BallsInfo) Size() int {
	sum := 4
	for _, v := range bsi.BallInfos {
		sum += v.Size()
	}
	return sum
}

// NewItems init BallInfos
func (bsi *BallsInfo) NewItems(length uint32) {
	bsi.BallInfos = make([]ball.Ball, length)
	for i := range bsi.BallInfos {
		bsi.BallInfos[i] = ball.NewBall()
	}
}

// Crop crop BallInfos
func (bsi *BallsInfo) Crop(length uint32) {
	bsi.BallInfos = bsi.BallInfos[:length]
}

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
	CollisionInfos []*CollisionInfo
}

// Length return length
func (csi *CollisionsInfo) Length() int {
	return len(csi.CollisionInfos)
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
}

// Crop crop CollisionInfos
func (csi *CollisionsInfo) Crop(length uint32) {
	csi.CollisionInfos = csi.CollisionInfos[:length]
}
