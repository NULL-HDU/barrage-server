package message

import (
	"barrage-server/ball"
	b "barrage-server/base"
	"barrage-server/libs/bufbo"
	"bytes"
	"testing"
)

func generateDisappearInfoBytes(num int) []byte {
	bs := make([]byte, 4+num*disappearInfoSize)
	bw := bufbo.NewBEBytesWriter(bs)

	// full ball id
	bw.PutUint32(uint32(num))

	for i := 0; i < num; i++ {
		bw.PutUint16(99)
	}

	return bs
}

// generateDisappearsInfo ...
func generateDisappearsInfo(num int) *DisappearsInfo {
	dsi := new(DisappearsInfo)

	dsi.IDs = make([]b.BallID, num)
	for i := range dsi.IDs {
		dsi.IDs[i] = 99
	}

	return dsi
}

// generateCollisionInfoBytes ...
func generateCollisionInfoBytes() []byte {
	bs := make([]byte, collisionInfoSize)
	bw := bufbo.NewBEBytesWriter(bs)

	// full ball id
	bw.PutUint32(uint32(uidA))
	bw.PutUint16(uint16(idA))
	bw.PutUint32(uint32(uidB))
	bw.PutUint16(uint16(idB))

	// damage
	bw.PutUint8(uint8(damageA))
	bw.PutUint8(uint8(damageB))

	// state
	bw.PutUint8(uint8(stateA))
	bw.PutUint8(uint8(stateB))

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
	csi := generateTestCollisionsInfo(int(length))

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

	ci := csi.CollisionInfos[3]
	uid, id, damage, state := ci.IDs[0].UID, ci.IDs[0].ID, ci.Damages[0], ci.States[0]
	if uid != b.UserID(uidA) || id != b.BallID(idA) ||
		damage != b.Damage(damageA) || state != ball.State(stateA) {
		t.Errorf("AInfo of items of CollisionsInfo is not correct! get uid: %v, id: %v, damage: %v, state: %d.", uid, id, damage, state)
	}
	uid, id, damage, state = ci.IDs[1].UID, ci.IDs[1].ID, ci.Damages[1], ci.States[1]
	if uid != b.UserID(uidB) || id != b.BallID(idB) ||
		damage != b.Damage(damageB) || state != ball.State(stateB) {
		t.Errorf("BInfo of items of CollisionsInfo is not correct! get uid: %v, id: %v, damage: %v, state: %d.", uid, id, damage, state)
	}

}

func TestBallsInfoMarshalListBinary(t *testing.T) {
	bsi := generateTestBallsInfo(4)

	bs, err := MarshalListBinary(bsi)
	if err != nil {
		t.Error(err)
	}
	t.Logf("bytes: % x", bs)

	if l1, l2 := bsi.Size(), len(bs); l1 != l2 {
		t.Errorf("Length of MarshalListBinary result should be %d, but get %d.", l1, l2)
	}
	if bs[3] != 4 {
		t.Errorf("Number of Balls should be %v, but get %d.", 4, bs[3])
	}
}

// TestBallsInfoUnmarshalListBinary ...
func TestBallsInfoUnmarshalListBinary(t *testing.T) {
	bsi := generateTestBallsInfo(40)
	bs, err := MarshalListBinary(bsi)
	if err != nil {
		t.Error(err)
	}

	batBsi := &BallsInfo{}
	n, _ := UnmarshalListBinary(batBsi, bs)
	if n != len(bs) {
		t.Errorf("Length of unmarshaled bytes should be %d, but get %d.", len(bs), n)
	}

	if l1, l2 := bsi.Length(), batBsi.Length(); l1 != l2 {
		t.Errorf("Length of BallsInfo should be %d, but get %d", l1, l2)
	}

}

func TestDisappearsInfoMarshallBinary(t *testing.T) {
	dsi := generateDisappearsInfo(20)
	bs, err := dsi.MarshalBinary()
	if err != nil {
		t.Error(err)
	}

	if size, rightSize := len(bs), 4+20*disappearInfoSize; size != rightSize {
		t.Errorf("Size of marshaled binary should be %d, but get %d.", rightSize, size)
	}

	dsi = new(DisappearsInfo)
	if err = dsi.UnmarshalBinary(bs); err != nil {
		t.Error(err)
	}
	if length, rightLen := len(dsi.IDs), 20; length != rightLen {
		t.Errorf("Length of items of marshaled binary should be %d, but get %d.", rightLen, length)
	}
	if v := dsi.IDs[0]; v != 99 {
		t.Errorf("Value of items of marshaled binary should be %d, but get %d.", 99, v)
	}
}
