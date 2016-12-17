package playground

import (
	"barrage-server/ball"
	b "barrage-server/base"
	m "barrage-server/message"
	tm "barrage-server/testLib/message"
	"testing"
)

// TestBallCache ...
func TestBallCache(t *testing.T) {
	bcmap := ballCache{}
	for i := 0; i < 20; i++ {
		bcmap[b.BallID(i)] = ball.NewBall()
	}

	bclist := bcmap.Balls()
	if length := len(bclist); length != 20 {
		t.Errorf("Result of Balls method is wrong, hope %d, get %d.", 20, length)
	}
}

// TODO: Error test
func TestPlayground(t *testing.T) {
	pg := NewPlayground().(*playground)

	// AddUser
	for i := 1; i <= 20; i++ {
		pg.AddUser(b.UserID(i))
	}
	if ucLen := len(pg.userCollisionCache); ucLen != 21 {
		t.Errorf("Length of user collisionInfo map is wrong, hope %d, get %d.", 21, ucLen)
	}
	if ubLen := len(pg.ballsGround); ubLen != 21 {
		t.Errorf("Length of user ballInfo map is wrong, hope %d, get %d.", 21, ubLen)
	}
	if ubsLen := len(pg.userBytesCache); ubsLen != 21 {
		t.Errorf("Length of user bytesInfo map is wrong, hope %d, get %d.", 21, ubsLen)
	}
	if unbcLen := len(pg.userNewBallsCache); unbcLen != 21 {
		t.Errorf("Length of user newBallInfo map is wrong, hope %d, get %d.", 21, unbcLen)
	}

	// PutPkg
	pi := tm.GenerateTestRandomPlaygroundInfo(1, 22, 20, 5, 15)

	if err := pg.PutPkg(pi); err != nil {
		t.Error(err)
	}
	if ubcLen := len(pg.ballsGround[1]); ubcLen != 0 {
		t.Errorf("Length of ballGround is wrong, hope %d, get %d.", 0, ubcLen)
	}
	if ucLen := len(pg.userCollisionCache[1]); ucLen != 5 {
		t.Errorf("Length of collisionCache is wrong, hope %d, get %d.", 5, ucLen)
	}
	if ubcLen := len(pg.userNewBallsCache[1]); ubcLen != 22 {
		t.Errorf("Length of newBallsCache is wrong, hope %d, get %d.", 22, ubcLen)
	}

	pg.PkgsForEachUser()
	if ubcLen := len(pg.ballsGround[1]); ubcLen != 22 {
		t.Errorf("Length of ballGround is wrong, hope %d, get %d.", 22, ubcLen)
	}

	// DeleteUser
	pg.DeleteUser(1)

	if ucLen := len(pg.userCollisionCache); ucLen != 20 {
		t.Errorf("Length of user collisionInfo map2is wrong, hope %d, get %d.", 20, ucLen)
	}
	if unbcLen := len(pg.userNewBallsCache); unbcLen != 20 {
		t.Errorf("Length of user newBallInfo map is wrong, hope %d, get %d.", 20, unbcLen)
	}
	if ubLen := len(pg.ballsGround); ubLen != 20 {
		t.Errorf("Length of user ballInfo map is wrong, hope %d, get %d.", 20, ubLen)
	}
	if ubsLen := len(pg.userBytesCache); ubsLen != 20 {
		t.Errorf("Length of user bytesInfo map is wrong, hope %d, get %d.", 20, ubsLen)
	}

	if uc0Len := len(pg.userCollisionCache[0]); uc0Len != 0 {
		t.Errorf("Length of collisionCache of Sys user is wrong, hope %d, get %d.", 0, uc0Len)
	}

	// PkgsForEachUser
	pi2 := tm.GenerateTestRandomPlaygroundInfo(2, 22, 20, 5, 0)
	if err := pg.PutPkg(pi2); err != nil {
		t.Error(err)
	}
	pg.PkgsForEachUser()
	pi2 = tm.GenerateTestRandomPlaygroundInfo(2, 22, 20, 5, 0)
	if err := pg.PutPkg(pi2); err != nil {
		t.Error(err)
	}

	pis := pg.PkgsForEachUser()
	if pisLen := len(pis); pisLen != 19 {
		t.Errorf("Length of playgroundInfo is wrong, hope %d, get %d.", 19, pisLen)
	}
	if ubcLen := len(pg.userBytesCache[2][newBallIndex].Buf); ubcLen != 0 {
		t.Errorf("Length of BytesCache is wrong, hope %d, get %d.", 0, ubcLen)
	}
	if ubcLen := len(pg.userBytesCache[2][ballsIndex].Buf); ubcLen != 0 {
		t.Errorf("Length of BytesCache is wrong, hope %d, get %d.", 0, ubcLen)
	}
	if ubcLen := len(pg.userBytesCache[2][collisionIndex].Buf); ubcLen != 0 {
		t.Errorf("Length of BytesCache is wrong, hope %d, get %d.", 0, ubcLen)
	}
	if ubcLen := len(pg.userBytesCache[2][bufferIndex].Buf); ubcLen != 0 {
		t.Errorf("Length of BytesCache is wrong, hope %d, get %d.", 0, ubcLen)
	}
	if ucLen := len(pg.userCollisionCache[2]); ucLen != 0 {
		t.Errorf("Length of collisionCache is wrong, hope %d, get %d.", 0, ucLen)
	}
	if ucLen := len(pg.userNewBallsCache[2]); ucLen != 0 {
		t.Errorf("Length of collisionCache is wrong, hope %d, get %d.", 0, ucLen)
	}

	if bgLen := len(pg.ballsGround[2]); bgLen != 44 {
		t.Errorf("Length of collisionCache is wrong, hope %d, get %d.", 44, bgLen)
	}

	pi = new(m.PlaygroundInfo)
	piForUnmarshal := pis[3]
	if piForUnmarshal.Receiver == 2 {
		piForUnmarshal = pis[4]
	}
	if err := pi.UnmarshalBinary(piForUnmarshal.CacheBytes); err != nil {
		t.Error(err)
	}
	if LenDisappear := len(pi.Disappears.IDs); LenDisappear != 0 {
		t.Errorf("Length of Disappears is wrong , hope %d, get %d.", 0, LenDisappear)
	}
	if LenNewBalls := len(pi.NewBalls.BallInfos); LenNewBalls != 22 {
		t.Errorf("Length of NewBallss is wrong , hope %d, get %d.", 22, LenNewBalls)
	}
	if LenDisplace := len(pi.Displacements.BallInfos); LenDisplace != 22 {
		t.Errorf("Length of LenDisplace is wrong , hope %d, get %d.", 22, LenDisplace)
	}
	if lencollision := len(pi.Collisions.CollisionInfos); lencollision != 5 {
		t.Errorf("length of lencollision is wrong , hope %d, get %d.", 5, lencollision)
	}

}
