package room

import (
	m "barrage-server/message"
	"encoding/binary"
	"testing"
)

// TestMergeInfoListBytes ...
func TestMergeInfoListBytes(t *testing.T) {
	var buffer []byte
	size := 0
	num := 0

	bsi := generateTestStruct(10)
	size = bsi.Size()
	num += 10

	bs, _ := m.MarshalListBinary(bsi)
	mergeInfoListBytes(&buffer, bs)

	if bufLen := len(buffer); bufLen != size {
		t.Errorf("Size of buffer is error, hope %d, get %d.", size, bufLen)
	}
	if bufNum := int(binary.BigEndian.Uint32(buffer)); bufNum != num {
		t.Errorf("Num of buffer is error, hope %d, get %d.", num, bufNum)
	}

	bsi = generateTestStruct(30)
	size += bsi.Size() - 4
	num += 30

	bs, _ = m.MarshalListBinary(bsi)
	mergeInfoListBytes(&buffer, bs)

	if bufLen := len(buffer); bufLen != size {
		t.Errorf("Size of buffer is error, hope %d, get %d.", size, bufLen)
	}
	if bufNum := int(binary.BigEndian.Uint32(buffer)); bufNum != num {
		t.Errorf("Num of buffer is error, hope %d, get %d.", num, bufNum)
	}

	bsi = generateTestStruct(60)
	size += bsi.Size() - 4
	num += 60

	bs, _ = m.MarshalListBinary(bsi)
	mergeInfoListBytes(&buffer, bs)

	if bufLen := len(buffer); bufLen != size {
		t.Errorf("Size of buffer is error, hope %d, get %d.", size, bufLen)
	}
	if bufNum := int(binary.BigEndian.Uint32(buffer)); bufNum != num {
		t.Errorf("Num of buffer is error, hope %d, get %d.", num, bufNum)
	}
}
