package room

import (
	b "barrage-server/base"
	m "barrage-server/message"
	"barrage-server/user"
	"sync"
)

// Hall is a struct for user who is not playing game in room,
// it holds all rooms and all online users, listens on connect info package from
// users then join user to aim room.
type Hall struct {
	rM      sync.RWMutex
	uM      sync.RWMutex
	statusM sync.RWMutex

	rooms map[b.RoomID]Room
	users map[b.UserID]user.User

	infoChan chan m.InfoPkg
	status   uint8
}

// NewHall create and init hall.
func NewHall() (h *Hall) {
	h = new(Hall)
	h.rooms = make(map[b.RoomID]Room)
	h.users = make(map[b.UserID]user.User)
	h.infoChan = make(chan m.InfoPkg, 10)

	return
}

// ID ...
func (h *Hall) ID() b.RoomID {
	return hallID
}

// UserJoin ...
func (h *Hall) UserJoin(u user.User) error {
	h.uM.Lock()
	defer h.uM.Unlock()

	_, ok := h.users[u.ID()]
	if !ok {
		h.users[u.ID()] = u
	}
	// aways rebind room of user.
	u.BindRoom(hallID, h.infoChan)

	return nil
}

// UserLeft ...
//
func (h *Hall) UserLeft(userID b.UserID) error {
	h.uM.Lock()
	defer h.uM.Unlock()

	delete(h.users, userID)
	return nil
}

// InfoChan ...
func (h *Hall) InfoChan() <-chan m.InfoPkg {
	return h.infoChan
}

// Status ...
func (h *Hall) Status() uint8 {
	h.statusM.RLock()
	defer h.statusM.RUnlock()

	return h.status
}

// HandleInfoPkg ...
func (h *Hall) HandleInfoPkg(ipkg m.InfoPkg) {
	switch t := ipkg.Type(); t {
	case m.InfoConnect:
		ci := ipkg.Body().(*m.ConnectInfo)
		logger.Infoln(ci.RID, ci.UID)
	case m.InfoDisconnect:
		di := ipkg.Body().(*m.DisconnectInfo)
		logger.Infoln(di.RID, di.UID)
	default:
		logger.Errorf("Invalid information package! type: %d.", t)
	}
}

// LoopOperation ...
func (h *Hall) LoopOperation() {

}

// CompareAndSetStatus ...
func (h *Hall) CompareAndSetStatus(oldStatus, newStatus uint8) (isSet bool) {
	h.statusM.Lock()
	defer h.statusM.Unlock()

	if h.status == oldStatus {
		h.status = newStatus
		isSet = true
	}
	return
}
