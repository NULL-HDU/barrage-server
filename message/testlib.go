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

	disappearID = 99
)

func generateBall() ball.Ball {
	b, _ := ball.NewUserAirplane(0, "airplane", 1, 0, 99, 99)
	return b
}

// generateTestBallsInfo ...
func generateTestBallsInfo(num int) *BallsInfo {
	bsi := &BallsInfo{}

	bsi.length = uint32(num)
	bsi.BallInfos = make([]ball.Ball, num)
	for i := 0; i < num; i++ {
		bsi.BallInfos[i] = generateBall()
	}

	return bsi
}

// generateTestCollisionsInfo ...
func generateTestCollisionsInfo(num int) *CollisionsInfo {
	csi := &CollisionsInfo{}
	csi.NewItems(uint32(num))
	fullBallIDA := b.FullBallID{
		UID: b.UserID(uidA),
		ID:  b.BallID(idA),
	}
	fullBallIDB := b.FullBallID{
		UID: b.UserID(uidB),
		ID:  b.BallID(idB),
	}

	for i := uint32(0); i < csi.length; i++ {
		csi.CollisionInfos[i] = &CollisionInfo{
			IDs:     []b.FullBallID{fullBallIDA, fullBallIDB},
			Damages: []b.Damage{b.Damage(damageA), b.Damage(damageB)},
			States:  []ball.State{stateA, stateB},
		}
	}

	return csi
}

// generateDisappearsInfo ...
func generateTestDisappearsInfo(num int) *DisappearsInfo {
	dsi := new(DisappearsInfo)

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
func GenerateTestPlaygroundInfo(sender b.UserID, niNum, diNum, ciNum, dsiNum int) *PlaygroundInfo {
	return &PlaygroundInfo{
		Sender:        sender,
		NewBalls:      generateTestBallsInfo(niNum),
		Displacements: generateTestBallsInfo(diNum),
		Collisions:    generateTestCollisionsInfo(ciNum),
		Disappears:    generateTestDisappearsInfo(dsiNum),
	}
}
