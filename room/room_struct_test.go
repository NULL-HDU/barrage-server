package room

import (
	b "barrage-server/base"
	m "barrage-server/message"
	tm "barrage-server/testLib/message"
	"testing"
	"time"
)

func init() {
	OpenGameHallAndRooms(b.OpenRoomIDs)
}

type testInfo struct {
	v1 byte
	v2 int
}

// Size ...
func (ti *testInfo) Size() int {
	return ti.v2
}

// return the []byte that value of each item is v1 and length is v2
func (ti *testInfo) MarshalBinary() ([]byte, error) {
	bs := make([]byte, ti.v2)

	bs[0] = ti.v1
	writedLen := 1
	for writedLen < ti.v2 {
		copy(bs[writedLen:], bs[:writedLen])
		writedLen *= 2
	}

	return bs, nil
}

func (ti *testInfo) UnmarshalBinary(bs []byte) error {
	if len(bs) == 0 {
		ti.v1, ti.v2 = 0, 0
	}

	ti.v1 = bs[0]
	ti.v2 = 1
	for _, vi := range bs[1:] {
		if vi == ti.v1 {
			ti.v2++
		} else {
			break
		}
	}

	return nil
}

type testInfoList struct {
	length   uint32
	infolist []testInfo
}

func (til *testInfoList) Length() int {
	return int(til.length)
}

// Size ...
func (til *testInfoList) Size() int {
	return 4 + int(til.length)*til.infolist[0].Size()
}

// Item ...
func (til *testInfoList) Item(index int) b.CommunicationData {
	return &til.infolist[index]
}

// NewItem ...
func (til *testInfoList) NewItems(length uint32) {
	til.infolist = make([]testInfo, length)
	til.length = length
}

// Crop ...
func (til *testInfoList) Crop(length uint32) {
	if til.length == length {
		return
	}
	til.infolist = til.infolist[:length]
	til.length = length
}

// generateTestStruct ...
func generateTestStruct(num int) *testInfoList {
	til := &testInfoList{}

	til.length = uint32(num)
	til.infolist = make([]testInfo, num)
	for i := 0; i < num; i++ {
		til.infolist[i] = testInfo{v1: 'b', v2: 10}
	}

	return til
}

type testUser struct {
	id  b.UserID
	rid b.RoomID

	infopkgChan chan<- m.InfoPkg
	checkFunc   func(bs []byte, itype m.InfoType)
}

// Play ...
func (tu *testUser) Play() error {
	return nil
}

// ID ...
func (tu *testUser) ID() b.UserID {
	return tu.id
}

// Room ...
func (tu *testUser) Room() b.RoomID {
	return tu.rid
}

func (tu *testUser) SendError(s string) {
	si := &m.SpecialMsgInfo{Message: s}
	tu.Send(si)
}

// Send ...
func (tu *testUser) Send(ipkg m.InfoPkg) {
	bs, err := ipkg.Body().MarshalBinary()
	itype := ipkg.Type()
	if err != nil {
		logger.Errorln(err)
	}

	if tu.checkFunc != nil {
		tu.checkFunc(bs, itype)
	} else {
		logger.Infoln("checkFunc is nil")
	}
}

// UploadInfo ...
func (tu *testUser) UploadInfo(infopkg m.InfoPkg) error {
	tu.infopkgChan <- infopkg
	return nil
}

// BindRoom ...
func (tu *testUser) BindRoom(id b.RoomID, c chan<- m.InfoPkg) {
	tu.rid = id
	tu.infopkgChan = c
}

// TestRoomUserJoinAndLeft ...
func TestRoomUserJoinAndLeftAndIDAndUsers(t *testing.T) {
	r := NewRoom(20)
	checkFunc := func(bs []byte, itype m.InfoType) {
		if itype != m.InfoConnected {
			t.Errorf("Receive wrong message type, hope %d, get %d.", m.InfoConnected, itype)
		}

		ci := new(m.ConnectedInfo)
		if err := ci.UnmarshalBinary(bs); err != nil {
			t.Error(err)
		}

		if ci.UID != 1 {
			t.Errorf("User id of ConnectedInfo should be %d, get %d.", 1, ci.UID)
		}
		if ci.RID != 20 {
			t.Errorf("Room id of ConnectedInfo should be %d, get %d.", 20, ci.RID)
		}
	}

	if id := r.ID(); id != 20 {
		t.Errorf("Room id is wrong, hope %d, get %d.", 20, id)
	}

	if status := r.Status(); status != roomClose {
		t.Errorf("Status of room should be %d, but get %d.", roomClose, status)
	}

	tu1 := &testUser{id: 1, checkFunc: checkFunc}
	tu2 := &testUser{id: 2}
	tu3 := &testUser{id: 3}
	tu4 := &testUser{id: 4}

	if err := r.UserJoin(tu1); err != nil {
		t.Error(err)
	}
	if err := r.UserJoin(tu2); err != nil {
		t.Error(err)
	}
	if err := r.UserJoin(tu3); err != nil {
		t.Error(err)
	}
	if err := r.UserJoin(tu4); err != nil {
		t.Error(err)
	}

	users := r.Users()
	if usersLen := len(users); usersLen != 4 {
		t.Errorf("Length of users is wrong, hope %d, get %d.", 4, usersLen)
	}

	if err := r.UserLeft(tu1.id); err != nil {
		t.Error(err)
	}
	if err := r.UserLeft(tu2.id); err != nil {
		t.Error(err)
	}

	if lenUser := len(commonHall.users); lenUser != 2 {
		t.Errorf("Number of users in CommonHall should be %d, but get %d.", 2, lenUser)
	}

	users = r.Users()
	if usersLen := len(users); usersLen != 2 {
		t.Errorf("Length of users is wrong, hope %d, get %d.", 2, usersLen)
	}

}

// TestRoomDisconnect ...
func TestRoomDisconnect(t *testing.T) {
	r := NewRoom(20)
	Open(r, time.Second)

	tu1 := &testUser{id: 1}
	if err := r.UserJoin(tu1); err != nil {
		t.Error(err)
	}

	if user := r.users[1]; user.ID() != tu1.id {
		t.Errorf("UserJoin error: user should has joined.")
	}

	di := &m.DisconnectInfo{UID: 1, RID: 20}
	tu1.UploadInfo(di)

	time.Sleep(time.Millisecond)
	if userLen := len(r.users); userLen != 0 {
		t.Errorf("Number of user should zero, but get %d.", userLen)
	}

	Close(r)
}

//TestRoomHandlePlaygroundInfo ...
func TestRoomHandlePlaygroundInfoAndPlaygroundBoardCast(t *testing.T) {
	pi1 := tm.GenerateTestRandomPlaygroundInfo(1, 30, 40, 15, 20)
	pi2 := tm.GenerateTestRandomPlaygroundInfo(2, 30, 40, 15, 20)
	pi3 := tm.GenerateTestRandomPlaygroundInfo(3, 30, 40, 15, 20)
	checkFunc := func(bs []byte, itype m.InfoType) {
		if itype != m.InfoPlayground {
			return
		}
		// Test playgroundBoardCast
		piBak := new(m.PlaygroundInfo)
		if err := piBak.UnmarshalBinary(bs); err != nil {
			t.Error(err)
		}

		if nbLen := piBak.NewBalls.Length(); nbLen != 0 {
			t.Errorf("Number of NewBalls is wrong, hope %d, get %d.", 0, nbLen)
		}
		if diLen := piBak.Displacements.Length(); diLen != 140 {
			t.Errorf("Number of Displacements is wrong, hope %d, get %d.", 140, diLen)
		}
		if ciLen := piBak.Collisions.Length(); ciLen != 30 {
			t.Errorf("Number of CollisionsInfo is wrong, hope %d, get %d.", 30, ciLen)
		}
		if dsiLen := len(piBak.Disappears.IDs); dsiLen != 0 {
			t.Errorf("Number of DisappearsInfo is wrong, hope %d, get %d.", 0, dsiLen)
		}

	}

	tu1 := &testUser{id: 1, checkFunc: checkFunc}
	tu2 := &testUser{id: 2, checkFunc: checkFunc}
	tu3 := &testUser{id: 3, checkFunc: checkFunc}

	r := NewRoom(20)
	Open(r, time.Second)

	if err := r.UserJoin(tu1); err != nil {
		t.Error(err)
	}
	if err := r.UserJoin(tu2); err != nil {
		t.Error(err)
	}
	if err := r.UserJoin(tu3); err != nil {
		t.Error(err)
	}

	tu1.UploadInfo(pi1)
	tu2.UploadInfo(pi2)
	tu3.UploadInfo(pi3)

	time.Sleep(time.Second)

	Close(r)
}
