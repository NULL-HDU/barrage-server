package message

import (
	"barrage-server/ball"
	b "barrage-server/base"
)

const (
	uidA = 10
	uidB = 100
	idA  = 1
	idB  = 1

	damageA = 120
	damageB = 130

	stateA = ball.Alive
	stateB = ball.Dead
)

func generateBall() ball.Ball {
	b, _ := ball.NewUserAirplane(0, "airplane", 1, 0, 99, 99)
	return b
}

// generateTestBallsInfo ...
func generateTestBallsInfo(num int) *BallsInfo {
	bsi := &BallsInfo{}

	bsi.length = uint32(num)
	bsi.ballInfos = make([]ball.Ball, num)
	for i := 0; i < num; i++ {
		bsi.ballInfos[i] = generateBall()
	}

	return bsi
}

// generateTestCollisionsInfo ...
func generateTestCollisionsInfo(num int) *CollisionsInfo {
	csi := &CollisionsInfo{}
	csi.NewItems(uint32(num))
	fullBallIDA := fullBallID{
		uid: b.UserID(uidA),
		id:  b.BallID(idA),
	}
	fullBallIDB := fullBallID{
		uid: b.UserID(uidB),
		id:  b.BallID(idB),
	}

	for i := uint32(0); i < csi.length; i++ {
		csi.CollisionInfos[i] = collisionInfo{
			ballIDs: []fullBallID{fullBallIDA, fullBallIDB},
			damages: []b.Damage{b.Damage(damageA), b.Damage(damageB)},
			states:  []ball.State{stateA, stateB},
		}
	}

	return csi
}

// GenerateTestPlaygroundInfo generate a testing playgroundInfo,
// ciNum is the number of collisionInfo in playgroundInfo
// diNum is the number of displacementInfo in playgroundInfo
// niNum is the number of newBallsInfo in playgroundInfo
func GenerateTestPlaygroundInfo(sender b.UserID, ciNum, diNum, niNum int) *PlaygroundInfo {
	return &PlaygroundInfo{
		Sender:        sender,
		Collisions:    generateTestCollisionsInfo(ciNum),
		Displacements: generateTestBallsInfo(diNum),
		NewBalls:      generateTestBallsInfo(niNum),
	}
}
