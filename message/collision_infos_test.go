package message

import (
	"barrage-server/ball"
	b "barrage-server/base"
	"barrage-server/libs/bufbo"
	"bytes"
	"testing"
)

const (
	uidA = 10
	uidB = 100
	idA  = 1
	idB  = 1

	damageA = 120
	damageB = 130

	stateA = ball.Alive
	stateB = ball.Dead
)

// generateTestCollisionsInfo ...
func generateTestCollisionsInfo(num uint8) *CollisionsInfo {
	csi := &CollisionsInfo{}
	csi.NewItems(uint32(num))
	fullBallIDA := fullBallID{
		uid: b.UserID(uidA),
		id:  b.BallID(idA),
	}
	fullBallIDB := fullBallID{
		uid: b.UserID(uidB),
		id:  b.BallID(idB),
	}

	for i := uint32(0); i < csi.length; i++ {
		csi.CollisionInfos[i] = collisionInfo{
			ballIDs: []fullBallID{fullBallIDA, fullBallIDB},
			damages: []b.Damage{b.Damage(damageA), b.Damage(damageB)},
			states:  []ball.State{stateA, stateB},
		}
	}

	return csi
}

// generateCollisionInfoBytes ...
func generateCollisionInfoBytes() []byte {
	bs := make([]byte, 28)
	bw := bufbo.NewBEBytesWriter(bs)

	// full ball id
	bw.PutUint64(uint64(uidA))
	bw.PutUint16(uint16(idA))
	bw.PutUint64(uint64(uidB))
	bw.PutUint16(uint16(idB))

	// damage
	bw.PutUint16(uint16(damageA))
	bw.PutUint16(uint16(damageB))

	// state
	isAlive, isKilled, _ := ball.AnalyseStateToBytes(stateA)
	bw.PutUint8(isAlive)
	bw.PutUint8(isKilled)
	isAlive, isKilled, _ = ball.AnalyseStateToBytes(stateB)
	bw.PutUint8(isAlive)
	bw.PutUint8(isKilled)

	return bs
}

// generateCollisionInfoBytes ...
func generateCollisionsInfoBytes(num uint8) []byte {
	var buffer bytes.Buffer

	buffer.Write([]byte{0, 0, 0, num})
	for i := uint8(0); i < num; i++ {
		buffer.Write(generateCollisionInfoBytes())
	}

	return buffer.Bytes()
}

// TestCollisionInfosMarsharlListBinary ...
func TestCollisionInfoMarsharlListBinary(t *testing.T) {
	length := uint8(9)
	csi := generateTestCollisionsInfo(length)

	bs, err := MarshalListBinary(csi)
	if err != nil {
		t.Error(err)
	}
	if bs2 := generateCollisionsInfoBytes(length); bytes.Compare(bs, bs2) != 0 {
		t.Errorf("MarshalListBinary result is not correct, hope %v, get %v.", bs2, bs)
	}
}

// TestUnmarshalListBinary ...
func TestCollisionsInfoUnmarshalListBinary(t *testing.T) {
	length := uint8(9)
	bs := generateCollisionsInfoBytes(length)

	csi := &CollisionsInfo{}
	n, err := UnmarshalListBinary(csi, bs)
	if err != nil {
		t.Error(err)
	}
	if n != len(bs) {
		t.Errorf("Length of unmarshaled bytes should be %d, but get %d.", n, len(bs))
	}

	uid, id, damage, state := csi.CollisionInfos[3].AInfo()
	if uid != b.UserID(uidA) || id != b.BallID(idA) ||
		damage != b.Damage(damageA) || state != ball.State(stateA) {
		t.Errorf("AInfo of items of CollisionsInfo is not correct! get uid: %v, id: %v, damage: %v, state: %d.", uid, id, damage, state)
	}
	uid, id, damage, state = csi.CollisionInfos[3].BInfo()
	if uid != b.UserID(uidB) || id != b.BallID(idB) ||
		damage != b.Damage(damageB) || state != ball.State(stateB) {
		t.Errorf("BInfo of items of CollisionsInfo is not correct! get uid: %v, id: %v, damage: %v, state: %d.", uid, id, damage, state)
	}

}
