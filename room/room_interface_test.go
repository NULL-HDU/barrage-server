package room

import (
	b "barrage-server/base"
	m "barrage-server/message"
	"testing"
	"time"
)

func init() {
	OpenGameHallAndRooms(b.OpenRoomIDs)
}

type testTiggler struct {
	id       b.RoomID
	status   uint8
	infoChan chan m.InfoPkg

	checkInfoPkg       func(ipkg m.InfoPkg)
	checkLoopOperation func()
}

func (tt *testTiggler) CompareAndSetStatus(oldStatus, newStatus uint8) (isSet bool) {
	if tt.status == oldStatus {
		tt.status = newStatus
		isSet = true
	}

	return
}

// ID ...
func (tt *testTiggler) ID() b.RoomID {
	return tt.id
}

// Status ...
func (tt *testTiggler) Status() uint8 {
	return tt.status
}

// HandleInfoPkg ...
func (tt *testTiggler) HandleInfoPkg(ipkg m.InfoPkg) {
	tt.checkInfoPkg(ipkg)
}

// InfoChan ...
func (tt *testTiggler) InfoChan() <-chan m.InfoPkg {
	return tt.infoChan
}

// LoopOperation ...
func (tt *testTiggler) LoopOperation() {
	tt.checkLoopOperation()
}

// TestTigglerOpenAndClose ...
func TestTigglerOpenAndClose(t *testing.T) {
	count := 0
	checkLoopOperation := func() {
		count += 1
	}
	checkInfoPkg := func(ipkg m.InfoPkg) {
		if infoType := ipkg.Type(); infoType != m.InfoGameOver {
			t.Errorf("Type of info package is error, hope %d, get %d.", m.InfoGameOver, infoType)
		}

	}

	testInfoChan := make(chan m.InfoPkg)
	tt := &testTiggler{
		id:                 99,
		infoChan:           testInfoChan,
		checkLoopOperation: checkLoopOperation,
		checkInfoPkg:       checkInfoPkg,
	}

	Open(tt, time.Second)
	if status := tt.Status(); status != roomOpen {
		t.Errorf("Status of open Tiggler should be %d, but get %d.", roomOpen, status)
	}

	time.Sleep(time.Second + 500*time.Millisecond)
	tt.infoChan <- &m.GameOverInfo{Overtype: 1}

	Close(tt)
	if status := tt.Status(); status != roomClose {
		t.Errorf("Status of open Tiggler should be %d, but get %d.", roomClose, status)
	}
	if count != 1 {
		t.Errorf("count is wrong, hope %d, get %d.", 1, count)
	}

	time.Sleep(time.Second)

}
