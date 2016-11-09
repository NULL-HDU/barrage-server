package playground

import (
	"barrage-server/ball"
	b "barrage-server/base"
	"math"
	"math/rand"
)

// CreateAirplaneInPlayGround create a airplane in random position.
func CreateAirplaneInPlayGround(c b.UserID, nickname string, role uint8, special uint16) (ball.Ball, error) {
	x, y := rand.Intn(math.MaxUint16), rand.Intn(math.MaxUint16)

	return ball.NewUserAirplane(c, nickname, role, special, uint16(x), uint16(y))
}
