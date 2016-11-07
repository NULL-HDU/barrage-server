package room

import (
	b "barrage-server/base"
	m "barrage-server/message"
	"barrage-server/user"
	"fmt"
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

// generateCacheMap create and init a cache map
func generateCacheMap() (cacheMap [][]byte) {
	cacheMap = make([][]byte, 4)
	return
}

// Room marshal and cache infoes from playground sorting them by info sender.
// When boardcast infoes from background, Room chooses and combines info bytes
// according to sender id.
type Room struct {
	mapM    sync.RWMutex
	statusM sync.RWMutex

	users map[b.UserID]user.User
	cache map[b.UserID][][]byte
	id    b.RoomID

	//TODO: add infoChan for playground
	infoChan chan m.InfoPkg

	// close: roomClose, open: roomOpen
	status uint8
}

// NewRoom create a room struct using room id.
func NewRoom(id b.RoomID) (r *Room) {
	r = new(Room)
	r.id = id
	r.infoChan = make(chan m.InfoPkg, 10)

	return
}

// ID ...
func (r *Room) ID() b.RoomID {
	return r.id
}

// Users ...
func (r *Room) Users() (users []b.UserID) {
	r.mapM.RLock()
	defer r.mapM.RUnlock()

	users = make([]b.UserID, 0, len(r.users))
	for k := range r.users {
		users = append(users, k)
	}

	return
}

// UserJoin ...
func (r *Room) UserJoin(u user.User) error {
	r.mapM.Lock()
	defer r.mapM.Unlock()

	_, ok := r.users[u.ID()]
	if !ok {
		r.users[u.ID()] = u
		// TODO: cache map also should be cache. (cache map list pool for every room.)
		r.cache[u.ID()] = generateCacheMap()
		u.BindRoom(r.id, r.infoChan)
	}

	return nil
}

// UserLeft ...
func (r *Room) UserLeft(userID b.UserID) error {

	r.mapM.Lock()
	defer r.mapM.Unlock()

	u, ok := r.users[userID]
	if !ok {
		return fmt.Errorf("User Left Error: Not found user %d.", userID)
	}

	JoinHall(u)
	delete(r.users, userID)
	delete(r.cache, userID)

	return nil
}

// handlePlayground ...
func (r *Room) handlePlayground(ipkg m.InfoPkg) {
	// TODO: add userID for infoPkg

}

// handleDisconnect ...
func (r *Room) handleDisconnect(ipkg m.InfoPkg) {
	dsi := ipkg.Body().(*m.DisconnectInfo)
	userID := dsi.UID
	roomID := dsi.RID

	if roomID != r.id {
		logger.Errorf("Disconnect Error: RoomID is wrong, hope %d, get %d.", r.id, roomID)
		return
	}

	r.UserLeft(userID)
}

// playgroundBoardCast ...
func (r *Room) playgroundBoardCast() {
}

// boardCast ...
func (r *Room) boardCast(info m.Info) {
	bs, err := info.MarshalBinary()
	if err != nil {
		logger.Errorln(err)
	}

	for _, v := range r.users {
		v.Send(bs)
	}
}

// InfoChan ...
func (r *Room) InfoChan() <-chan m.InfoPkg {
	return r.infoChan
}

// LoopOperation wrap playgroundBoardCast.
func (r *Room) LoopOperation() {
	r.playgroundBoardCast()
}

// HandleInfoPkg ...
func (r *Room) HandleInfoPkg(ipkg m.InfoPkg) {
	switch t := ipkg.Type(); t {
	case m.InfoPlayground:
		r.handlePlayground(ipkg)
	case m.InfoDisconnect:
		r.handleDisconnect(ipkg)

	// flowing two type is unusable now.
	case m.InfoAirplaneCreated:
	case m.InfoSpecialMessage:
	default:
		logger.Errorf("Invalid information package! type: %d.", t)
	}
}

// Status check whether room is close.
//
// Before user sending infopkg into infoChan, they should call this fucntion to check
// whether room is close.
func (r *Room) Status() uint8 {
	r.statusM.RLock()
	defer r.statusM.RUnlock()

	return r.status
}

// CompareAndSetStatus ...
func (r *Room) CompareAndSetStatus(oldStatus, newStatus uint8) (isSet bool) {
	r.statusM.Lock()
	defer r.statusM.Unlock()

	if r.status == oldStatus {
		r.status = newStatus
		isSet = true
	}

	return
}
