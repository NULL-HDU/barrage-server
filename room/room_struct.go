package room

import (
	b "barrage-server/base"
	m "barrage-server/message"
	pg "barrage-server/playground"
	"barrage-server/user"
	"encoding/binary"
	"fmt"
	"sync"
)

var rmLimit = b.RoomMembersLimit

const (
	collisionIndex = iota
	displaceIndex
	newballIndex

	// cache data for send to self client.
	bufferIndex
)

// generateCacheMap create and init a cache map
func generateCacheMap() (cacheMap []roomCache) {
	cacheMap = make([]roomCache, 4)
	return
}

type roomCache struct {
	Num uint32
	Buf []byte
}

// clearCache set Num to 0 and truncate the Buf
func clearCache(cache *roomCache) {
	cache.Num = 0
	cache.Buf = cache.Buf[:0]
}

// cacheInfoListBytes add InfoList bytes into roomCache.
// This function will count the number of Info and connect rest bytes.
func cacheInfoListBytes(rc *roomCache, nb []byte) {
	if len(nb) <= 4 {
		return
	}

	// calculate new number
	n1 := rc.Num
	n2 := binary.BigEndian.Uint32(nb)
	num := n1 + n2
	if n1 > num {
		// n1 > num,  there are too many informations, drop them
		logger.Warnf("Too many informations in a list! Get %d.", n2)
		return
	}

	// set new number
	rc.Num = num
	// connect them
	rc.Buf = append(rc.Buf, nb[4:]...)
}

// Room marshal and cache infoes from playground sorting them by info sender.
// When boardcast infoes from background, Room chooses and combines info bytes
// according to sender id.
type Room struct {
	mapM    sync.RWMutex
	statusM sync.RWMutex

	users map[b.UserID]user.User
	cache map[b.UserID][]roomCache
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
	r.users = make(map[b.UserID]user.User)
	r.cache = make(map[b.UserID][]roomCache)
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
func (r *Room) UserJoin(u user.User, name string) error {
	r.mapM.Lock()
	defer r.mapM.Unlock()

	if len(r.users) >= rmLimit {
		return errRoomIsFull
	}

	_, ok := r.users[u.ID()]
	if ok {
		return errUserAlreadyJoin
	}

	uid := u.ID()
	r.users[uid] = u
	// TODO: cache map also should be cache. (cache map list pool for every room.)
	newCacheMap := generateCacheMap()
	r.cache[uid] = newCacheMap
	u.BindRoom(r.id, r.infoChan)

	airplane, err := pg.CreateAirplaneInPlayGround(uid, name, 1, 0)
	if err != nil {
		return err
	}

	// cache airplane message
	bs, err := airplane.MarshalBinary()
	if err != nil {
		return err
	}
	newCacheMap[newballIndex].Num = 1
	newCacheMap[newballIndex].Buf = bs

	// send airplaneCreatedInfo
	aci := new(m.AirplaneCreatedInfo)
	aci.Airplane = airplane
	u.Send(aci)

	logger.Infof("User %d join room %d. \n", u.ID(), r.id)

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

	logger.Infof("User %d left room %d. \n", userID, r.id)

	return nil
}

// handlePlayground add playgroundInfo data into the cache of pi.Sender in room
func (r *Room) handlePlayground(pi *m.PlaygroundInfo) {
	r.mapM.Lock()
	defer r.mapM.Unlock()

	byteCache, ok := r.cache[pi.Sender]
	if !ok {
		logger.Errorf("Not find user %d in room cache map %d. \n", pi.Sender, r.id)
		return
	}

	bs, err := m.MarshalListBinary(pi.Collisions)
	if err != nil {
		logger.Errorln(err)
	}
	cacheInfoListBytes(&byteCache[collisionIndex], bs)

	bs, err = m.MarshalListBinary(pi.Displacements)
	if err != nil {
		logger.Errorln(err)
	}
	cacheInfoListBytes(&byteCache[displaceIndex], bs)

	bs, err = m.MarshalListBinary(pi.NewBalls)
	if err != nil {
		logger.Errorln(err)
	}
	cacheInfoListBytes(&byteCache[newballIndex], bs)
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

// constructApartBytesFor append bytes of partIndex in r.cache of other user.
func (r *Room) constructApartBytesFor(uid b.UserID, partIndex int) {
	bufferCache := &r.cache[uid][bufferIndex]
	lenOffset := len(bufferCache.Buf)
	listItemCount := uint32(0)

	bufferCache.Buf = append(bufferCache.Buf, make([]byte, 4)...)

	for _uid, rc := range r.cache {
		if _uid == uid {
			continue
		}
		listItemCount += rc[partIndex].Num
		bufferCache.Buf = append(bufferCache.Buf, rc[partIndex].Buf...)
	}

	binary.BigEndian.PutUint32(bufferCache.Buf[lenOffset:], listItemCount)
}

// constructBytesFor connect bytes slice into playgroundInfoPkg.
func (r *Room) constructBytesFor(uid b.UserID) {
	r.constructApartBytesFor(uid, collisionIndex)
	r.constructApartBytesFor(uid, displaceIndex)
	r.constructApartBytesFor(uid, newballIndex)
}

// playgroundBoardCast ...
func (r *Room) playgroundBoardCast() {
	r.mapM.Lock()
	defer r.mapM.Unlock()

	// construct message and send
	for uid, user := range r.users {
		cache, ok := r.cache[uid]
		if !ok {
			logger.Errorf("Not find user %d in room cache map %d.", uid, r.id)
			continue
		}
		r.constructBytesFor(uid)

		// while lenght of bs is greater than 12, playground info is not empty.
		if bs := cache[bufferIndex].Buf; len(bs) > 12 {
			// TODO: write bytes_cache for every info pkg.
			pi := new(m.PlaygroundInfo)
			if err := pi.UnmarshalBinary(bs); err != nil {
				logger.Errorf("PlaygroundInfo Create Error: %s\n", err)
			}

			// send bytes in a new goroutine
			user.Send(pi)
		}
	}

	// clear cache
	for _, cache := range r.cache {
		clearCache(&cache[collisionIndex])
		clearCache(&cache[displaceIndex])
		clearCache(&cache[newballIndex])
		clearCache(&cache[bufferIndex])
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
