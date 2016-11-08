package room

import (
	"testing"
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
