package message

import (
	b "barrage-server/base"
	"bytes"
	"encoding/binary"
	"errors"
	"testing"
)

type testInfo struct {
	v1 byte
	v2 int
}

// Size ...
func (ti *testInfo) Size() int {
	return ti.v2
}

// return the []byte that value of each item is v1 and length is v2
func (ti *testInfo) MarshalBinary() ([]byte, error) {
	if ti.v2 == -1 {
		return nil, errors.New("error for testing.")
	}

	bs := make([]byte, ti.v2)

	bs[0] = ti.v1
	writedLen := 1
	for writedLen < ti.v2 {
		copy(bs[writedLen:], bs[:writedLen])
		writedLen *= 2
	}

	return bs, nil
}

func (ti *testInfo) UnmarshalBinary(bs []byte) error {
	if len(bs) == 0 {
		ti.v1, ti.v2 = 0, 0
	}

	ti.v1 = bs[0]
	ti.v2 = 1
	for _, vi := range bs[1:] {
		if vi == ti.v1 {
			ti.v2++
		} else {
			break
		}
	}

	return nil
}

type testInfoList struct {
	length   uint32
	infolist []testInfo
}

func (til *testInfoList) Length() int {
	return int(til.length)
}

// Item ...
func (til *testInfoList) Item(index int) b.CommunicationData {
	return &til.infolist[index]
}

// NewItem ...
func (til *testInfoList) NewItems(length uint32) {
	til.infolist = make([]testInfo, length)
	til.length = length
}

// Crop ...
func (til *testInfoList) Crop(length uint32) {
	if til.length == length {
		return
	}
	til.infolist = til.infolist[:length]
	til.length = length
}

// generateTestBytes ...
func generateTestBytes() []byte {
	var buffer bytes.Buffer

	buffer.Write([]byte{0, 0, 0, 4})
	buffer.Write(bytes.Repeat([]byte{'a'}, 10))
	buffer.Write(bytes.Repeat([]byte{'b'}, 10))
	buffer.Write(bytes.Repeat([]byte{'c'}, 9))
	buffer.Write(bytes.Repeat([]byte{'d'}, 9))

	return buffer.Bytes()
}

// generateTestStruct ...
func generateTestErrorStruct() *testInfoList {
	til := &testInfoList{}

	til.length = 4
	til.infolist = make([]testInfo, 4)
	til.infolist[0] = testInfo{v1: 'a', v2: -1}
	til.infolist[1] = testInfo{v1: 'b', v2: -1}
	til.infolist[2] = testInfo{v1: 'c', v2: -1}
	til.infolist[3] = testInfo{v1: 'd', v2: -1}

	return til
}

// generateTestStruct ...
func generateTestStruct() *testInfoList {
	til := &testInfoList{}

	til.length = 4
	til.infolist = make([]testInfo, 4)
	til.infolist[0] = testInfo{v1: 'a', v2: 10}
	til.infolist[1] = testInfo{v1: 'b', v2: 10}
	til.infolist[2] = testInfo{v1: 'c', v2: 9}
	til.infolist[3] = testInfo{v1: 'd', v2: 9}

	return til
}

func TestUnmarshalListBinary(t *testing.T) {
	til := &testInfoList{}
	bs := generateTestBytes()
	t.Logf("bytes: % x", bs)

	n, err := UnmarshalListBinary(til, bs)
	if err != nil {
		t.Error(err)
	}
	if n != 42 {
		t.Errorf("n should be 42, but get %d.", n)
	}

	if til.length != 4 {
		t.Errorf("length of testInfoList should be 4, but get %d.", til.length)
	}
	if v1 := til.infolist[0].v1; v1 != 'a' {
		t.Errorf("value of v1 of first element testInfoList should be %d, but get %d", 'a', v1)
	}
	if v2 := til.infolist[0].v2; v2 != 10 {
		t.Errorf("value of v2 of first element testInfoList should be 10, but get %d", v2)
	}
	if v1 := til.infolist[1].v1; v1 != 'b' {
		t.Errorf("value of v1 of second element testInfoList should be %d, but get %d", 'b', v1)
	}
	if v2 := til.infolist[1].v2; v2 != 10 {
		t.Errorf("value of v2 of second element testInfoList should be 10, but get %d", v2)
	}
	if v1 := til.infolist[2].v1; v1 != 'c' {
		t.Errorf("value of v1 of third element testInfoList should be %d, but get %d", 'c', v1)
	}
	if v2 := til.infolist[2].v2; v2 != 9 {
		t.Errorf("value of v2 of third element testInfoList should be 9, but get %d", v2)
	}
	if v1 := til.infolist[3].v1; v1 != 'd' {
		t.Errorf("value of v1 of forth element testInfoList should be %d, but get %d", 'd', v1)
	}
	if v2 := til.infolist[3].v2; v2 != 9 {
		t.Errorf("value of v2 of forth element testInfoList should be 9, but get %d", v2)
	}

	bs[3] = 0
	n, err = UnmarshalListBinary(til, bs)
	if err != ErrEmptyInfo {
		t.Errorf("UnmarshalListBinary should return %v, but get %v.", ErrEmptyInfo, err)
	}

}

// TestMarshalListBinary ...
func TestMarshalListBinary(t *testing.T) {
	bs, err := MarshalListBinary(generateTestStruct())
	if err != nil {
		t.Error(err)
	}
	if bs2 := generateTestBytes(); bytes.Compare(bs2, bs) != 0 {
		t.Errorf("MarshalListBinary result is not correct, hope %v, get %v.", bs2, bs)
	}
}

// TestMarshalListBinary ...
func TestMarshalListBinaryError(t *testing.T) {
	errStr := generateTestErrorStruct()
	bs, err := MarshalListBinary(errStr)
	if err != nil {
		t.Error(err)
	}
	if bsLen := len(bs); bsLen != 4 {
		t.Errorf("MarshalListBinary result is not correct, hope 4, get %v", bsLen)
	}
	if v := int(binary.BigEndian.Uint32(bs)); v != 0 {
		t.Errorf("Head value of MarshalListBinary result is not correct, hope 0, get %v", v)
	}

}
