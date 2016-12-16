package message

import (
	"barrage-server/ball"
	b "barrage-server/base"
	m "barrage-server/message"
	tball "barrage-server/testLib/ball"
	"math/rand"
)

const (
	UidA = 10
	UidB = 100
	IdA  = 1
	IdB  = 1

	DamageA = 120
	DamageB = 130

	StateA = ball.Alive
	StateB = ball.Dead

	DisappearID = 99
)

var existBall = map[b.UserID][]ball.Ball{}

func generateBall() ball.Ball {
	b := ball.NewBallWithSpecialID(0, 99)
	return b
}

// generateTestBallsInfo ...
func generateTestBallsInfo(num int) *m.BallsInfo {
	bsi := &m.BallsInfo{}
	bsi.BallInfos = make([]ball.Ball, num)
	for i := 0; i < num; i++ {
		bsi.BallInfos[i] = generateBall()
	}

	return bsi
}

// generateTestCollisionsInfo ...
func generateTestCollisionsInfo(num int) *m.CollisionsInfo {
	csi := &m.CollisionsInfo{}
	csi.NewItems(uint32(num))
	fullBallIDA := b.FullBallID{
		UID: b.UserID(UidA),
		ID:  b.BallID(IdA),
	}
	fullBallIDB := b.FullBallID{
		UID: b.UserID(UidB),
		ID:  b.BallID(IdB),
	}

	for i := 0; i < csi.Length(); i++ {
		csi.CollisionInfos[i] = &m.CollisionInfo{
			IDs:     []b.FullBallID{fullBallIDA, fullBallIDB},
			Damages: []b.Damage{b.Damage(DamageA), b.Damage(DamageB)},
			States:  []ball.State{StateA, StateB},
		}
	}

	return csi
}

// generateDisappearsInfo ...
func generateTestDisappearsInfo(num int) *m.DisappearsInfo {
	dsi := new(m.DisappearsInfo)

	dsi.IDs = make([]b.BallID, num)
	for i := range dsi.IDs {
		dsi.IDs[i] = 99
	}

	return dsi
}

// GenerateTestPlaygroundInfo generate a testing playgroundInfo,
// ciNum is the number of CollisionInfo in playgroundInfo
// diNum is the number of displacementInfo in playgroundInfo
// niNum is the number of newBallsInfo in playgroundInfo
func GenerateTestPlaygroundInfo(sender b.UserID, niNum, diNum, ciNum, dsiNum int) *m.PlaygroundInfo {
	return &m.PlaygroundInfo{
		Sender:        sender,
		NewBalls:      generateTestBallsInfo(niNum),
		Displacements: generateTestBallsInfo(diNum),
		Collisions:    generateTestCollisionsInfo(ciNum),
		Disappears:    generateTestDisappearsInfo(dsiNum),
	}
}

// GenerateCollisionsInfoFromBall ...
func GenerateCollisionsInfoFromBalls(balls []ball.Ball) *m.CollisionsInfo {
	csi := &m.CollisionsInfo{}
	csi.NewItems(uint32(len(balls)))
	// random A id
	fullBallIDA := b.FullBallID{
		UID: b.UserID(rand.Intn(300)),
		ID:  b.BallID(rand.Intn(600)),
	}
	for i := 0; i < csi.Length(); i++ {
		csi.CollisionInfos[i] = &m.CollisionInfo{
			IDs: []b.FullBallID{
				fullBallIDA,
				b.FullBallID{
					UID: b.UserID(balls[i].UID()),
					ID:  b.BallID(balls[i].ID()),
				},
			},
			Damages: []b.Damage{b.Damage(DamageA), b.Damage(DamageB)},
			States:  []ball.State{ball.Alive, ball.Dead},
		}
	}

	return csi
}

// GenerateDisappearsInfoFromBalls ...
func GenerateDisappearsInfoFromBalls(balls []ball.Ball) *m.DisappearsInfo {
	dsi := make([]b.BallID, len(balls))

	for i := range balls {
		dsi[i] = balls[i].ID()
	}

	return &m.DisappearsInfo{
		IDs: dsi,
	}
}

func GenerateTestRandomPlaygroundInfo(sender b.UserID, niNum, diNum, ciNum, dsiNum int) *m.PlaygroundInfo {
	newBalls := tball.GenerateRandomIDBall(sender, ciNum+dsiNum+niNum)
	var userBalls []ball.Ball
	if v, ok := existBall[sender]; ok {
		userBalls = v[:diNum]
	}
	collisionBall := GenerateCollisionsInfoFromBalls(newBalls[:ciNum])
	disappearBall := GenerateDisappearsInfoFromBalls(newBalls[ciNum : ciNum+dsiNum])
	return &m.PlaygroundInfo{
		Sender: sender,
		NewBalls: &m.BallsInfo{
			BallInfos: newBalls,
		},
		Displacements: &m.BallsInfo{
			BallInfos: userBalls,
		},
		Collisions: collisionBall,
		Disappears: disappearBall,
	}
}
