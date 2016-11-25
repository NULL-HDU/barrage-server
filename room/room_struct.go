package room

import (
	b "barrage-server/base"
	m "barrage-server/message"
	pg "barrage-server/playground"
	"barrage-server/user"
	"sync"
)

var rmLimit = b.RoomMembersLimit

// Room marshal and cache infoes from playground sorting them by info sender.
// When boardcast infoes from background, Room chooses and combines info bytes
// according to sender id.
type Room struct {
	mapM    sync.RWMutex
	statusM sync.RWMutex

	users      map[b.UserID]user.User
	playground pg.Playground
	id         b.RoomID

	//TODO: add infoChan for playground
	infoChan chan m.InfoPkg

	// close: roomClose, open: roomOpen
	status uint8
}

// NewRoom create a room struct using room id.
func NewRoom(id b.RoomID) (r *Room) {
	r = new(Room)
	r.id = id
	r.users = make(map[b.UserID]user.User)
	r.playground = pg.NewPlayground()
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
// now this function will create an Airplane.
func (r *Room) UserJoin(u user.User) error {
	if err := r.userJoin(u); err != nil {
		return err
	}

	uid := u.ID()

	// send connected info back to front end.
	ci := &m.ConnectedInfo{UID: uid, RID: r.id}
	u.Send(ci)

	logger.Infof("User %d join room %d. \n", uid, r.id)

	return nil
}

// userJoin ...
func (r *Room) userJoin(u user.User) error {
	r.mapM.Lock()
	defer r.mapM.Unlock()

	if len(r.users) >= rmLimit {
		return errRoomIsFull
	}

	uid := u.ID()
	_, ok := r.users[uid]
	if ok {
		return errUserAlreadyJoin
	}

	r.users[uid] = u
	r.playground.AddUser(uid)
	u.BindRoom(r.id, r.infoChan)

	return nil
}

// UserLeft ...
func (r *Room) UserLeft(userID b.UserID) error {

	if err := r.userLeft(userID); err != nil {
		return err
	}

	logger.Infof("User %d left room %d. \n", userID, r.id)
	return nil
}

// userLeft ...
func (r *Room) userLeft(userID b.UserID) error {
	r.mapM.Lock()
	defer r.mapM.Unlock()

	u, ok := r.users[userID]
	if !ok {
		return errUserNotFound
	}

	JoinHall(u)
	r.playground.DeleteUser(userID)
	delete(r.users, userID)

	return nil
}

// handlePlayground add playgroundInfo data into the cache of pi.Sender in room
func (r *Room) handlePlayground(pi *m.PlaygroundInfo) {
	if err := r.playground.PutPkg(pi); err != nil {
		if err == pg.ErrNotFoundUser {
			logger.Errorf("Not find user %d in room cache map %d. \n", pi.Sender, r.id)
		} else {
			logger.Errorln(err)
		}
	}

}

// handleDisconnect ...
func (r *Room) handleDisconnect(dsi *m.DisconnectInfo) {
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
	pis := r.playground.PkgsForEachUser()

	r.mapM.Lock()
	defer r.mapM.Unlock()

	for _, pi := range pis {
		u, ok := r.users[pi.Receiver]
		if !ok {
			logger.Errorf("Not find user %d in room cache map %d.", u.ID(), r.id)
			continue
		}
		if len(pi.CacheBytes) > 16 {
			u.Send(pi)
		}
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
	var err string

	switch t := ipkg.Type(); t {
	case m.InfoPlayground:
		pi, ok := ipkg.Body().(*m.PlaygroundInfo)
		if !ok {
			err = "InfoPkg fails to be convert into PlaygroundInfo."
			break
		}
		r.handlePlayground(pi)
	case m.InfoDisconnect:
		dsi, ok := ipkg.Body().(*m.DisconnectInfo)
		if !ok {
			err = "InfoPkg fails to be convert into DisconnectInfo."
			break
		}
		r.handleDisconnect(dsi)

	// flowing two type is unusable now.
	case m.InfoAirplaneCreated:
	case m.InfoSpecialMessage:
	default:
		logger.Errorf("Invalid information package! type: %d.\n", t)
	}

	if err != "" {
		logger.Errorln(err)
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
