package room

import (
	b "barrage-server/base"
	m "barrage-server/message"
	"testing"
	"time"
)

// TestHallUserJoinAndLeft ...
func TestHallUserJoinAndLeftAndID(t *testing.T) {
	r := NewHall()
	if id := r.ID(); id != 0 {
		t.Errorf("Room id is wrong, hope %d, get %d.", hallID, id)
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

	users := r.users
	if usersLen := len(users); usersLen != 4 {
		t.Errorf("Length of users is wrong, hope %d, get %d.", 4, usersLen)
	}

	if err := r.UserLeft(tu1.id); err != nil {
		t.Error(err)
	}
	if err := r.UserLeft(tu2.id); err != nil {
		t.Error(err)
	}

	users = r.users
	if usersLen := len(users); usersLen != 2 {
		t.Errorf("Length of users is wrong, hope %d, get %d.", 2, usersLen)
	}

}

// TestHallHandleInfoPkg ...
func TestHallHandleInfoPkg(t *testing.T) {
	r := NewHall()
	count := 0
	checkFunc := func(bs []byte, itype m.InfoType) {
		if itype != m.InfoSpecialMessage {
			count = -1
			return
		}
		count++
	}

	tu1 := &testUser{
		id:        1,
		checkFunc: checkFunc,
	}

	if err := r.UserJoin(tu1); err != nil {
		t.Error(err)
	}

	ci := new(m.ConnectInfo)
	ci.UID = 1
	ci.RID = 20

	Open(r, time.Second)

	// not exist
	tu1.UploadInfo(ci)
	time.Sleep(10 * time.Millisecond)
	if count != 1 {
		t.Errorf("tu1 should join failed and receive a special message, but get count %d.", count)
	}
	count = 0

	r.rooms[20] = NewRoom(20)
	tu1.UploadInfo(ci)
	time.Sleep(100 * time.Millisecond)
	if count != -1 {
		t.Errorf("tu1 should join success and set count to -1, but get count %d.", count)
	}
	count = 0

	// full
	for i := 0; i < rmLimit+1; i++ {
		tu := &testUser{
			id:        b.UserID(i + 2),
			checkFunc: checkFunc,
		}
		if err := r.UserJoin(tu); err != nil {
			t.Error(err)
		}
		ci := new(m.ConnectInfo)
		ci.UID = b.UserID(i + 2)
		ci.RID = 20
		tu.UploadInfo(ci)
	}
	time.Sleep(10 * time.Millisecond)
	if count != 1 {
		t.Errorf("tu should join failed and set count to 1, but get count %d.", count)
	}
	count = 0

	// rejoin
	ci = new(m.ConnectInfo)
	ci.UID = 1
	ci.RID = 20
	tu := &testUser{
		id:        23,
		checkFunc: checkFunc,
	}
	if err := r.UserJoin(tu); err != nil {
		t.Error(err)
	}
	tu.UploadInfo(ci)
	time.Sleep(10 * time.Millisecond)
	if count != 1 {
		t.Errorf("tu1 should join failed and set count to 1, but get count %d.", count)
	}
	count = 0

}
