package room

import (
	b "barrage-server/base"
	m "barrage-server/message"
	"testing"
)

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
	checkFunc   func(bs []byte)
}

// ID ...
func (tu *testUser) ID() b.UserID {
	return tu.id
}

// Room ...
func (tu *testUser) Room() b.RoomID {
	return tu.rid
}

// Send ...
func (tu *testUser) Send(bs []byte) {
	tu.checkFunc(bs)
}

// UploadInfo ...
func (tu *testUser) UploadInfo(infopkg m.InfoPkg) {
	tu.infopkgChan <- infopkg
}

// BindRoom ...
func (tu *testUser) BindRoom(id b.RoomID, c chan<- m.InfoPkg) {
	tu.rid = id
	tu.infopkgChan = c
}

// TestRoomUserJoinAndLeft ...
func TestRoomUserJoinAndLeftAndIDAndUsers(t *testing.T) {
	r := NewRoom(20)
	if id := r.ID(); id != 20 {
		t.Errorf("Room id is wrong, hope %d, get %d.", 20, id)
	}

	if status := r.Status(); status != roomClose {
		t.Errorf("Status of room should be %d, but get %d.", roomClose, status)
	}

	tu1 := &testUser{id: 1}
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

	users = r.Users()
	if usersLen := len(users); usersLen != 2 {
		t.Errorf("Length of users is wrong, hope %d, get %d.", 2, usersLen)
	}

}

// TestMergeInfoListBytes ...
func TestCacheInfoListBytes(t *testing.T) {
	var buffer roomCache
	size := 0
	num := 0

	bsi := generateTestStruct(10)
	size = bsi.Size() - 4
	num += 10

	bs, _ := m.MarshalListBinary(bsi)
	cacheInfoListBytes(&buffer, bs)

	if bufLen := len(buffer.Buf); bufLen != size {
		t.Errorf("Size of buffer is error, hope %d, get %d.", size, bufLen)
	}
	if bufNum := int(buffer.Num); bufNum != num {
		t.Errorf("Num of buffer is error, hope %d, get %d.", num, bufNum)
	}

	bsi = generateTestStruct(30)
	size += bsi.Size() - 4
	num += 30

	bs, _ = m.MarshalListBinary(bsi)
	cacheInfoListBytes(&buffer, bs)

	if bufLen := len(buffer.Buf); bufLen != size {
		t.Errorf("Size of buffer is error, hope %d, get %d.", size, bufLen)
	}
	if bufNum := int(buffer.Num); bufNum != num {
		t.Errorf("Num of buffer is error, hope %d, get %d.", num, bufNum)
	}

	bsi = generateTestStruct(60)
	size += bsi.Size() - 4
	num += 60

	bs, _ = m.MarshalListBinary(bsi)
	cacheInfoListBytes(&buffer, bs)

	if bufLen := len(buffer.Buf); bufLen != size {
		t.Errorf("Size of buffer is error, hope %d, get %d.", size, bufLen)
	}
	if bufNum := int(buffer.Num); bufNum != num {
		t.Errorf("Num of buffer is error, hope %d, get %d.", num, bufNum)
	}
}
