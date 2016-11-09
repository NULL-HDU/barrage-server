package message

import (
	"testing"
)

func TestBallsInfoMarshalListBinary(t *testing.T) {
	bsi := generateTestBallsInfo(4)

	bs, err := MarshalListBinary(bsi)
	if err != nil {
		t.Error(err)
	}
	t.Logf("bytes: % x", bs)

	if l1, l2 := bsi.Size(), len(bs); l1 != l2 {
		t.Errorf("Length of MarshalListBinary result should be %d, but get %d.", l1, l2)
	}
	if bs[3] != 4 {
		t.Errorf("Number of Balls should be %v, but get %d.", 4, bs[3])
	}
}

// TestBallsInfoUnmarshalListBinary ...
func TestBallsInfoUnmarshalListBinary(t *testing.T) {
	bsi := generateTestBallsInfo(40)
	bs, err := MarshalListBinary(bsi)
	if err != nil {
		t.Error(err)
	}

	batBsi := &BallsInfo{}
	n, _ := UnmarshalListBinary(batBsi, bs)
	if n != len(bs) {
		t.Errorf("Length of unmarshaled bytes should be %d, but get %d.", len(bs), n)
	}

	if l1, l2 := bsi.Length(), batBsi.Length(); l1 != l2 {
		t.Errorf("Length of BallsInfo should be %d, but get %d", l1, l2)
	}

}
