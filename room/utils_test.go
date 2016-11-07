package room

import (
	m "barrage-server/message"
	"encoding/binary"
	"testing"
)

// generateTestBallsInfo ...
func generateTestBallsInfo(num int) *m.BallsInfo {
	bsi := &m.BallsInfo{}
	bsi.NewItems(uint32(num))
	return bsi
}

// TestMergeInfoListBytes ...
func TestMergeInfoListBytes(t *testing.T) {
	var buffer []byte

	bsi1 := generateTestBallsInfo(10)
	bs, _ := m.MarshalListBinary(bsi1)
	mergeInfoListBytes(&buffer, bs)
}
