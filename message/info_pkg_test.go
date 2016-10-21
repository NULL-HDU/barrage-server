package message

import (
	"barrage-server/ball"
	b "barrage-server/base"
	"testing"
)

// TestGameOverInfo
func TestGameOverInfo(t *testing.T) {

	// MarshalBinary
	goi := &GameOverInfo{uint8(1)}
	bs, err := goi.MarshalBinary()
	if err != nil {
		t.Error(err)
	}

	if len(bs) != 1 {
		t.Errorf("Length of marshaled bytes should be 1, but get %d.", len(bs))
	}
	if bs[0] != uint8(1) {
		t.Errorf("Value of marshaled bytes should be { 1 }, but get %v.", bs)
	}

	// UnmarshalBinary
	bs = []byte{1}
	goi = &GameOverInfo{}
	err = goi.UnmarshalBinary(bs)
	if err != nil {
		t.Error(err)
	}

	if ot := goi.Overtype; ot != uint8(1) {
		t.Errorf("Value of Overtype should be 1, but get %v.", ot)
	}
}

// TestSpecialMsgInfo ...
func TestSpecialMsgInfo(t *testing.T) {
	testStr := "Testing information"

	// MarshalBinary
	smi := &SpecialMsgInfo{testStr}
	bs, err := smi.MarshalBinary()
	if err != nil {
		t.Error(err)
	}

	if smi.Size() != 20 {
		t.Errorf("Size of Marshaled bytes should be 20, but get %d.", smi.Size())
	}

	if l1, l2 := len(bs), smi.Size(); l1 != l2 {
		t.Errorf("Result of Marshaled bytes is not correct, hope %d, get %d.", l2, l1)
	}
	if strFromByte := string(bs[1:]); strFromByte != testStr {
		t.Errorf("Value of Marshalend bytes is not correct, hope %s, get %s.", testStr, strFromByte)
	}

	// UnmarshalBinary
	bs = append([]byte{19}, []byte(testStr)...)
	smi = &SpecialMsgInfo{}
	err = smi.UnmarshalBinary(bs)
	if err != nil {
		t.Error(err)
	}

	if smi.Size() != 20 {
		t.Errorf("Size of Marshaled bytes should be 20, but get %d.", smi.Size())
	}
	if smi.Message != testStr {
		t.Errorf("Message of smi should be %s, but get %s.", testStr, smi.Message)
	}

}

// TestAirplaneCreatedInfo ...
func TestAirplaneCreatedInfo(t *testing.T) {
	// MarshalBinary
	airplane, err := ball.NewUserAirplane(0, "Tester", 1, 2, 99, 99)
	aci := &AirplaneCreatedInfo{airplane}
	bs, err := aci.MarshalBinary()
	if err != nil {
		t.Error(err)
	}

	if aci.Size() != 32 {
		t.Errorf("Size of Marshaled bytes should be 32, but get %d.", aci.Size())
	}
	if l1, l2 := len(bs), aci.Size(); l1 != l2 {
		t.Errorf("Result of Marshaled bytes is not correct, hope %d, get %d.", l2, l1)
	}

	// UnmarshalBinary
	aci = &AirplaneCreatedInfo{}
	err = aci.UnmarshalBinary(bs)
	if uid := aci.Airplane.UID(); uid != b.UserID(0) {
		t.Errorf("User Id of Unmarshaled AirplaneCreatedInfo should be %v, but get %v.", b.UserID(0), uid)
	}
}

// TestDisconnectInfo ...
func TestDisconnectInfo(t *testing.T) {
	// MarshalBinary
	di := &DisconnectInfo{b.UserID(2333), b.RoomID(1)}
	bs, err := di.MarshalBinary()
	if err != nil {
		t.Error(err)
	}

	if l1, l2 := len(bs), di.Size(); l1 != l2 {
		t.Errorf("Result of Marshaled bytes is not correct, hope %d, get %d.", l2, l1)
	}

	// UnmarshalBinary
	di = &DisconnectInfo{}
	err = di.UnmarshalBinary(bs)
	if uid := di.UID; uid != b.UserID(2333) {
		t.Errorf("User Id of Unmarshaled DisconnectInfo should be %v, but get %v.", b.UserID(2333), uid)
	}
	if rid := di.RID; rid != b.RoomID(1) {
		t.Errorf("Room Id of Unmarshaled DisconnectInfo should be %v, but get %v.", b.RoomID(1), rid)
	}
}

// TestConnectInfo ...
func TestConnectInfo(t *testing.T) {
	// MarshalBinary
	ci := &ConnectInfo{b.UserID(666666), "Tester", b.RoomID(1), 1}
	bs, err := ci.MarshalBinary()
	if err != nil {
		t.Error(err)
	}

	if l1, l2 := len(bs), ci.Size(); l1 != l2 {
		t.Errorf("Result of Marshaled bytes is not correct, hope %d, get %d.", l2, l1)
	}

	// UnmarshalBinary
	ci = &ConnectInfo{}
	err = ci.UnmarshalBinary(bs)
	if uid := ci.UID; uid != b.UserID(666666) {
		t.Errorf("User Id of Unmarshaled ConnectInfo should be %v, but get %v.", b.UserID(666666), uid)
	}
	if nickname := ci.Nickname; nickname != "Tester" {
		t.Errorf("Nickname of Unmarshaled ConnectInfo should be %v, but get %v.", "Tester", nickname)
	}
	if rid := ci.RID; rid != b.RoomID(1) {
		t.Errorf("Room Id of Unmarshaled ConnectInfo should be %v, but get %v.", b.RoomID(1), rid)
	}
	if troop := ci.Troop; troop != uint8(1) {
		t.Errorf("Troop of Unmarshaled ConnectInfo should be %v, but get %v.", uint8(1), troop)
	}
}

// TestPlaygroundInfo ...
func TestPlaygroundInfo(t *testing.T) {
	// MarshalBinary
	pi := &PlaygroundInfo{
		Collisions:    generateTestCollisionsInfo(9),
		Displacements: generateTestBallsInfo(30),
		NewBalls:      generateTestBallsInfo(20),
	}
	bs, err := pi.MarshalBinary()
	if err != nil {
		t.Error(err)
	}
	if l1, l2 := len(bs), pi.Size(); l1 != l2 {
		t.Errorf("Result of Marshaled bytes is not correct, hope %d, get %d.", l2, l1)
	}

	// UnmarshalBinary
	pi = &PlaygroundInfo{}
	err = pi.UnmarshalBinary(bs)
	if err != nil {
		t.Error(err)
	}
	if csiLen := pi.Collisions.Length(); csiLen != 9 {
		t.Errorf("Length of PlaygroundInfo Collisions should be %d, but get %d.", 9, csiLen)
	}
	if diLen := pi.Displacements.Length(); diLen != 30 {
		t.Errorf("Length of PlaygroundInfo Collisions should be %d, but get %d.", 30, diLen)
	}
	if niLen := pi.NewBalls.Length(); niLen != 20 {
		t.Errorf("Length of PlaygroundInfo Collisions should be %d, but get %d.", 20, niLen)
	}
}
