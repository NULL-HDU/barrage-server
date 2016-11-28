package ball

import (
	"barrage-server/ball"
	b "barrage-server/base"
	"barrage-server/libs/bufbo"
	"math/rand"
)

type TestBall struct {
	Uid      b.UserID
	Id       b.BallID
	Nickname string
}

func (bl *TestBall) UID() b.UserID {
	return bl.Uid
}

func (bl *TestBall) ID() b.BallID {
	return bl.Id
}

func (bl *TestBall) Size() int {
	return 7 + len(bl.Nickname)
}

func (bl *TestBall) MarshalBinary() ([]byte, error) {
	bs := make([]byte, bl.Size())
	bw := bufbo.NewBEBytesWriter(bs)

	bw.PutUint32(uint32(bl.Uid))
	bw.PutUint16(uint16(bl.Id))

	nicknameLen := len(bl.Nickname)
	bw.PutUint8(uint8(nicknameLen))
	bw.PutStr(bl.Nickname)

	return bs, nil
}

func (bl *TestBall) UnmarshalBinary(data []byte) error {
	br := bufbo.NewBEBytesReader(data)

	bl.Uid = b.UserID(br.Uint32())
	bl.Id = b.BallID(br.Uint16())
	bl.Nickname = br.Str(int(br.Uint8()))
	return nil
}

// GenerateRandomIDBall ...
func GenerateRandomIDBall(uid b.UserID, num int) []ball.Ball {
	balls := make([]ball.Ball, num)
	for i := range balls {
		balls[i] = ball.NewBallWithSpecialID(uid, b.BallID(rand.Int()))
	}

	return balls
}
