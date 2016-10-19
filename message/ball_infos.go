package message

import (
	"barrage-server/ball"
	b "barrage-server/base"
)

// BallsInfo is used for ball informations transimission.
type BallsInfo struct {
	length    uint32
	ballInfos []ball.Ball
}

// Length return length
func (bsi *BallsInfo) Length() int {
	return int(bsi.length)
}

// Item return item of ballInfos.
func (bsi *BallsInfo) Item(index int) b.CommunicationData {
	return bsi.ballInfos[index]
}

// Size return the number of bytes after marshed
func (bsi *BallsInfo) Size() int {
	sum := 4
	for _, v := range bsi.ballInfos {
		sum += v.Size()
	}
	return sum
}

// NewItems init ballInfos
func (bsi *BallsInfo) NewItems(length uint32) {
	bsi.ballInfos = make([]ball.Ball, length)
	for i := range bsi.ballInfos {
		bsi.ballInfos[i] = ball.NewBall()
	}
	bsi.length = length
}

// Crop crop ballInfos
func (bsi *BallsInfo) Crop(length uint32) {
	if bsi.length == length {
		return
	}
	bsi.ballInfos = bsi.ballInfos[:length]
	bsi.length = length
}
