package playground

import (
	"barrage-server/ball"
	b "barrage-server/base"
	m "barrage-server/message"
	"errors"
	"sync"
)

var (
	errNotFoundUser = errors.New("Not found user.")
)

type ballCache struct {
	m map[b.BallID]ball.Ball
	//cache result of Balls()
	balls       []ball.Ball
	didMChanged bool
}

// Len ...
func (bc *ballCache) Len() int {
	return len(bc.m)
}

// Get ...
func (bc *ballCache) Get(id b.BallID) (result ball.Ball, ok bool) {
	result, ok = bc.m[id]
	return
}

func (bc *ballCache) Set(ba ball.Ball) {
	bc.m[ba.ID()] = ba
	bc.didMChanged = true
}

// Delete ...
func (bc *ballCache) Delete(id b.BallID) {
	delete(bc.m, id)
	bc.didMChanged = true
}

// Balls ...
func (bc *ballCache) Balls() []ball.Ball {
	if !bc.didMChanged {
		return bc.balls
	}

	bc.didMChanged = false
	bc.balls = make([]ball.Ball, 0, len(bc.m))
	for _, v := range bc.m {
		bc.balls = append(bc.balls, v)
	}

	return bc.balls
}

// Playground cache and pack up the collisionInfo, displacementInfo and newBallsInfo,
// it keep a user-ball map which be synchronous according to displacementInfo. it cache
// collisionInfo for sending to frontend.
type Playground interface {
	// init Playground
	Init()
	// Add user by uid.
	AddUser(b.UserID)
	// delete user by uid.
	DeleteUser(b.UserID)
	// construct playgroundInfo for every user in playground, then
	// map users by cb with parameters uid and its playgroundInfo
	// package.
	PkgForEachUser(cb func(b.UserID, m.PlaygroundInfo))
	// cache and pack up the infos in playgroundInfo.
	PutPkg(pi m.PlaygroundInfo) error
}

type playground struct {
	mapM sync.RWMutex

	userCollision map[b.UserID][]*m.CollisionInfo
	userBallCache map[b.UserID]*ballCache
}

// NewPlayground create default implement of Playground.
func NewPlayground() Playground {
	pg := &playground{
		userCollision: make(map[b.UserID][]*m.CollisionInfo),
		userBallCache: make(map[b.UserID]*ballCache),
	}

	pg.AddUser(b.SysID)
	return pg
}

// PutPkg ...
func (pg *playground) PutPkg(pi m.PlaygroundInfo) error {

}

// packUpPkg ...
func (pg *playground) packUpPkgs(pi m.PlaygroundInfo) error {
	pg.mapM.Lock()
	defer pg.mapM.Unlock()

	uid := pi.Sender
	_, ok := pg.userBallCache[uid]
	if !ok {
		return errNotFoundUser
	}

	validCollisionInfos = make([]*m.CollisionInfo, 0, pi.Collisions.Length())
	pg.userCollision[uid] = append(pg.userCollision[uid], pi.Collisions.CollisionInfos...)

}

// AddUser ...
func (pg *playground) AddUser(uid b.UserID) {
	pg.mapM.Lock()
	defer pg.mapM.Unlock()

	_, ok := pg.userBallCache[uid]
	if ok {
		return
	}

	bc := new(ballCache)
	bc.m = make(map[b.BallID]ball.Ball)
	pg.userCollision[uid] = make([]*m.CollisionInfo, 0)
	pg.userBallCache[uid] = bc
}

// DeleteUser ...
func (pg *playground) DeleteUser(uid b.UserID) {
	pg.mapM.Lock()
	defer pg.mapM.Unlock()

	_, ok := pg.userBallCache[uid]
	if !ok {
		return
	}

	// move collisionInfos of the uid user to SysID user.
	pg.userCollision[b.SysID] = append(pg.userCollision[b.SysID], pg.userCollision[uid]...)

	// change all ball to be collisionInfo and add then to SysID user.
	// TODO: maybe balls of delete user or death user could be food or block
	newCollisions := make([]*m.CollisionInfo, pg.userBallCache[uid].Len())
	fid1 := b.FullBallID{UID: b.SysID, ID: b.SysID}
	for i, ub := range pg.userBallCache[uid].Balls() {
		fid2 := b.FullBallID{UID: uid, ID: ub.ID()}
		newCollisions[i] = &m.CollisionInfo{
			IDs:     []b.FullBallID{fid1, fid2},
			Damages: []b.Damage{0, 0},
			States:  []ball.State{ball.Alive, ball.Disappear},
		}
	}
	pg.userCollision[b.SysID] = append(pg.userCollision[b.SysID], newCollisions...)

	delete(pg.userCollision, uid)
	delete(pg.userBallCache, uid)
}
