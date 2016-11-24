package playground

import (
	"barrage-server/ball"
	b "barrage-server/base"
	"math/rand"
)

// CreateAirplaneInPlayGround create a airplane in random position.
func CreateAirplaneInPlayGround(c b.UserID, nickname string, role uint8, special uint16) (ball.Ball, error) {
	x, y := rand.Intn(b.PlayGroundWidth), rand.Intn(b.PlayGroundHeight)

	return ball.NewUserAirplane(c, nickname, role, special, uint16(x), uint16(y))
}
