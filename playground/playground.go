package playground

import (
	"barrage-server/ball"
	b "barrage-server/base"
	m "barrage-server/message"
	"encoding/binary"
	"errors"
	"sync"
)

var logger = b.Log

const (
	newBallIndex = iota
	ballsIndex
	collisionIndex

	// cache data for send to self client.
	bufferIndex
)

var (
	// ErrNotFoundUser ...
	ErrNotFoundUser = errors.New("Not found user.")
)

type ballCache map[b.BallID]ball.Ball

// Balls ...
func (bc ballCache) Balls() (balls []ball.Ball) {
	balls = make([]ball.Ball, 0, len(bc))
	for _, v := range bc {
		balls = append(balls, v)
	}

	return balls
}

type bytesCache struct {
	Num uint32
	Buf []byte
}

// generateCacheMap create and init a cache map
func generateCacheMap() (cacheMap []bytesCache) {
	cacheMap = make([]bytesCache, 4)
	return
}

// clearCache set Num to 0 and truncate the Buf
func clearCache(cache *bytesCache) {
	cache.Num = 0
	cache.Buf = cache.Buf[:0]
}

// Playground cache and pack up the collisionInfo, displacementInfo and newBallsInfo,
// it keep a user-ball map which be synchronous according to displacementInfo. it cache
// collisionInfo for sending to frontend.
type Playground interface {
	// Add user by uid.
	AddUser(b.UserID)
	// delete user by uid.
	DeleteUser(b.UserID)
	// construct playgroundInfo for every user in playground, then
	// map users by cb with parameters uid and its playgroundInfo
	// package.
	PkgsForEachUser() []*m.PlaygroundInfo
	// cache and pack up the infos in playgroundInfo.
	PutPkg(pi *m.PlaygroundInfo) error
}

type playground struct {
	mapM sync.RWMutex

	// all balls in playground
	ballsGround map[b.UserID]ballCache
	// collect new balls templately
	userNewBallsCache map[b.UserID]ballCache
	// collect collisionInfo templately
	userCollisionCache map[b.UserID][]*m.CollisionInfo

	// not concurrent secrity. only be used by fillPlaygroundInfo.
	userBytesCache map[b.UserID][]bytesCache
}

// NewPlayground create default implement of Playground.
func NewPlayground() Playground {
	pg := &playground{
		userCollisionCache: make(map[b.UserID][]*m.CollisionInfo),
		ballsGround:        make(map[b.UserID]ballCache),
		userNewBallsCache:  make(map[b.UserID]ballCache),
		userBytesCache:     make(map[b.UserID][]bytesCache),
	}

	pg.AddUser(b.SysID)
	return pg
}

// preCompileForEachUser compile Balls and CollisionInfos in the ballsGround and userCollisionCache of
// a user, and put the compiled bytes into userBytesCache, the next operating will take advantage of
// these bytes.
func (pg *playground) preCompileForEachUser() {
	for uid := range pg.ballsGround {
		bsc := pg.userBytesCache[uid]

		// compile and cache collisionInfo
		csi := new(m.CollisionsInfo)
		csi.CollisionInfos = pg.userCollisionCache[uid]
		bs, err := m.MarshalListBinary(csi)
		if err != nil {
			logger.Errorln(err)
		}

		bsc[collisionIndex].Num = uint32(csi.Length())
		bsc[collisionIndex].Buf = append(bsc[collisionIndex].Buf, bs[4:]...)

		// compile and cache ballsIndex
		bi := new(m.BallsInfo)
		bi.BallInfos = pg.ballsGround[uid].Balls()
		bs, err = m.MarshalListBinary(bi)
		if err != nil {
			logger.Errorln(err)
		}

		bsc[ballsIndex].Num = uint32(bi.Length())
		bsc[ballsIndex].Buf = append(bsc[ballsIndex].Buf, bs[4:]...)

		// compile and cache newBallIndex
		nbi := new(m.BallsInfo)
		nbi.BallInfos = pg.userNewBallsCache[uid].Balls()
		bs, err = m.MarshalListBinary(nbi)
		if err != nil {
			logger.Errorln(err)
		}

		bsc[newBallIndex].Num = uint32(nbi.Length())
		bsc[newBallIndex].Buf = append(bsc[newBallIndex].Buf, bs[4:]...)
	}
}

// fillPlaygroundInfo construct a playgroundInfo, it append all compiled infos to CacheBytes of
// the playgroundInfo, but other attributes of playgroundInfo is empty. So this playgroundInfo should
// be only used for send to user without other operating.
func (pg *playground) fillPlaygroundInfo(uid b.UserID, pi *m.PlaygroundInfo) {
	pi.Receiver = uid
	bufferCache := &pg.userBytesCache[uid][bufferIndex]

	// append newBallsInfo
	pg.constructApartBytesFor(uid, newBallIndex)
	// append displacementInfo
	pg.constructApartBytesFor(uid, ballsIndex)
	// append collisionInfo
	pg.constructApartBytesFor(uid, collisionIndex)
	// append disappearInfos
	bufferCache.Buf = append(bufferCache.Buf, []byte{0, 0, 0, 0}...)

	pi.CacheBytes = append(pi.CacheBytes, bufferCache.Buf...)
}

// constructApartBytesFor append bytes of partIndex in userBytesCache of other user.
func (pg *playground) constructApartBytesFor(uid b.UserID, partIndex int) {
	bufferCache := &pg.userBytesCache[uid][bufferIndex]
	lenOffset := len(bufferCache.Buf)
	listItemCount := uint32(0)

	bufferCache.Buf = append(bufferCache.Buf, []byte{0, 0, 0, 0}...)

	for k, bsc := range pg.userBytesCache {
		if k == uid {
			continue
		}
		if bsc[partIndex].Num != 0 {
			listItemCount += bsc[partIndex].Num
			bufferCache.Buf = append(bufferCache.Buf, bsc[partIndex].Buf...)
		}
	}

	binary.BigEndian.PutUint32(bufferCache.Buf[lenOffset:], listItemCount)
}

func (pg *playground) cleanCacheForEachUser() {
	for uid := range pg.ballsGround {
		pg.userCollisionCache[uid] = pg.userCollisionCache[uid][:0]
		bsc := pg.userBytesCache[uid]

		clearCache(&bsc[newBallIndex])
		clearCache(&bsc[ballsIndex])
		clearCache(&bsc[collisionIndex])
		clearCache(&bsc[bufferIndex])

		for k, v := range pg.userNewBallsCache[uid] {
			pg.ballsGround[uid][k] = v
		}
		pg.userNewBallsCache[uid] = ballCache{}
	}
}

// PkgsForEachUser ...
func (pg *playground) PkgsForEachUser() (pis []*m.PlaygroundInfo) {
	pg.mapM.RLock()
	defer pg.mapM.RUnlock()

	// pre-compile and cache result
	pg.preCompileForEachUser()

	// construct playgroundInfo for each user.
	// not include Sys user
	pis = make([]*m.PlaygroundInfo, len(pg.ballsGround)-1)
	for i := range pis {
		pis[i] = new(m.PlaygroundInfo)
	}

	count := 0
	for uid := range pg.ballsGround {
		if uid == b.SysID {
			continue
		}
		pg.fillPlaygroundInfo(uid, pis[count])
		count++
	}

	// clean cache
	pg.cleanCacheForEachUser()

	return
}

// PutPkg ...
func (pg *playground) PutPkg(pi *m.PlaygroundInfo) error {
	return pg.packUpPkgs(pi)
}

// packUpPkg put Balls of newBallsInfo into ballsGround of the Sender, set Balls of DisplacementsInfo
// into ballsGround of ther Sender, cache collisionInfo of CollisionsInfo in userCollisionCache of the Sender,
// and delete Balls of Disappears of the Sender.
func (pg *playground) packUpPkgs(pi *m.PlaygroundInfo) error {
	pg.mapM.Lock()
	defer pg.mapM.Unlock()

	uid := pi.Sender
	_, ok := pg.ballsGround[uid]
	if !ok {
		return ErrNotFoundUser
	}

	// TODO: more check
	bg := pg.ballsGround[uid]
	nb := pg.userNewBallsCache[uid]

	// newBallsInfo, add new ball to ballCache map of uid
	for _, v := range pi.NewBalls.BallInfos {
		nb[v.ID()] = v
	}

	// displacementInfo, if ball is not in ballsGround, the ball should in newBallsCache
	// then modify existing balls in the appropriate place.
	for _, v := range pi.Displacements.BallInfos {
		if _, ok := nb[v.ID()]; ok {
			nb[v.ID()] = v
			continue
		}
		if _, ok := bg[v.ID()]; ok {
			bg[v.ID()] = v
			continue
		}
	}

	// collisionInfo, base check and cache them to userCollisionCache
	validCollisionInfos := make([]*m.CollisionInfo, 0, pi.Collisions.Length())
	for _, v := range pi.Collisions.CollisionInfos {
		// there are an unknow error
		if len(v.States) == 0 {
			logger.Errorln("Invalid collisionInfo: ", v)
			continue
		}
		if v.States[1] != ball.Alive {
			if deleted := pg.checkAndDeleteBall(v.IDs[1].UID, v.IDs[1].ID); !deleted {
				continue
			}
		}
		if v.States[0] != ball.Alive {
			if deleted := pg.checkAndDeleteBall(v.IDs[0].UID, v.IDs[0].ID); !deleted {
				continue
			}
		}

		validCollisionInfos = append(validCollisionInfos, v)
	}

	pg.userCollisionCache[uid] = append(pg.userCollisionCache[uid], validCollisionInfos...)

	// disappearInfos
	for _, v := range pi.Disappears.IDs {
		delete(bg, v)
		delete(nb, v)
	}

	return nil
}

// checkAndDeleteBall, a valid collisionInfo means if one of collision ball is death or disappear,
// this ball should be found in ballsGround, otherwise the collisionInfo is invalid.
func (pg *playground) checkAndDeleteBall(uid b.UserID, id b.BallID) (deleted bool) {
	v, ok := pg.ballsGround[uid]
	if !ok {
		return false
	}

	_, ok = v[id]
	if !ok {
		v = pg.userNewBallsCache[uid]
		_, ok = v[id]
		if !ok {
			return false
		}
	}

	delete(v, id)
	return true
}

// AddUser ...
func (pg *playground) AddUser(uid b.UserID) {
	pg.mapM.Lock()
	defer pg.mapM.Unlock()

	_, ok := pg.ballsGround[uid]
	if ok {
		return
	}

	pg.userCollisionCache[uid] = make([]*m.CollisionInfo, 0)
	pg.ballsGround[uid] = ballCache{}
	pg.userNewBallsCache[uid] = ballCache{}
	pg.userBytesCache[uid] = generateCacheMap()
}

// changeBallsToCollisionInfoAndPutToSysCache ...
func (pg *playground) changeBallsToCollisionInfoAndPutToSysCache(uid b.UserID, bc ballCache) {
	// change all ball to be collisionInfo and add then to SysID user.
	// TODO: maybe balls of delete user or death user could be food or block
	newCollisions := make([]*m.CollisionInfo, len(bc))
	fid1 := b.FullBallID{UID: b.SysID, ID: b.SysID}
	for i, ub := range bc.Balls() {
		fid2 := b.FullBallID{UID: uid, ID: ub.ID()}
		newCollisions[i] = &m.CollisionInfo{
			IDs:     []b.FullBallID{fid1, fid2},
			Damages: []b.Damage{0, 0},
			States:  []ball.State{ball.Alive, ball.Disappear},
		}
	}
	pg.userCollisionCache[b.SysID] = append(pg.userCollisionCache[b.SysID], newCollisions...)
}

// DeleteUser ...
func (pg *playground) DeleteUser(uid b.UserID) {
	pg.mapM.Lock()
	defer pg.mapM.Unlock()

	_, ok := pg.ballsGround[uid]
	if !ok {
		return
	}

	// move collisionInfos of the uid user to SysID user.
	pg.userCollisionCache[b.SysID] = append(pg.userCollisionCache[b.SysID], pg.userCollisionCache[uid]...)

	// change and move balls in ballsGround into userCollisionCache of SysID user;
	// pg.changeBallsToCollisionInfoAndPutToSysCache(uid, pg.ballsGround[uid])
	// change and move balls in userNewBallsCache into userCollisionCache of SysID user;
	// pg.changeBallsToCollisionInfoAndPutToSysCache(uid, pg.userNewBallsCache[uid])

	delete(pg.userCollisionCache, uid)
	delete(pg.ballsGround, uid)
	delete(pg.userNewBallsCache, uid)
	// TODO: cache map also should be cache. (cache map list pool for every room.)
	delete(pg.userBytesCache, uid)
}
