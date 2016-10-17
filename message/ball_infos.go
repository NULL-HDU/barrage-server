package message

import (
	"barrage-server/ball"
	b "barrage-server/base"
)

const (
	lengthOfBallInfo = 41
)

// BallsInfo is used for ball informations transimission.
type BallsInfo struct {
	length    uint32
	ballInfos []ball.Ball
}

// Length return length
func (bsi *BallsInfo) Length() uint32 {
	return bsi.length
}

// SizeOfItem return number of bytes of collisionInfos.
func (bsi *BallsInfo) SizeOfItem() int {
	return lengthOfBallInfo
}

// Item return item of ballInfos.
func (bsi *BallsInfo) Item(index int) b.CommunicationData {
	return bsi.ballInfos[index]
}

// NewItems init ballInfos
func (bsi *BallsInfo) NewItems(length uint32) {
	bsi.ballInfos = make([]ball.Ball, length)
	bsi.length = length
}

// Crop crop ballInfos
func (bsi *BallsInfo) Crop(length uint32) {
	bsi.ballInfos = bsi.ballInfos[:length]
	bsi.length = length
}
