package ball

import (
	"barrage-server/ball"
	b "barrage-server/base"
	"math/rand"
)

// GenerateRandomIDBall ...
func GenerateRandomIDBall(uid b.UserID, num int) []ball.Ball {
	balls := make([]ball.Ball, num)
	for i := range balls {
		balls[i] = ball.NewBallWithSpecialID(uid, b.BallID(rand.Intn(100000)))
	}

	return balls
}
